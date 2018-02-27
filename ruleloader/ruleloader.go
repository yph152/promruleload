package ruleloader

import (
	"encoding/json"
	"fmt"
	"github.com/yph152/goproject/promeloadrules/cache"
	"github.com/yph152/goproject/promeloadrules/config"
	"github.com/yph152/goproject/promeloadrules/etcd_client"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Rule struct {
	Namespace    string `json:"namespace,omitempty"`   //业务名
	Server_name  string `json:"server_name,omitempty"` //服务名字
	Alertname    string `json:"alertname,omitempty"`   //告警名称
	Condition    string `json:"condition,omitempty"`   //告警条件
	Time         string `json:"time,omitempty"`        //延迟告警时间
	Value        string `json:"value,omitempty"`
	Level        string `json:"level,omitempty"`        //告警等级
	Trigername   string `json:"trigername,omitempty"`   //触发器名称
	Hostgroup    string `json:"hostgroup,omitempty"`    //用户组
	Templatename string `json:"templatename,omitempty"` //告警模板名称
	Summary      string `json:"summary,omitempty"`
	Description  string `json:"description,omitempty"` //告警内容
}

type RuleLoader struct {
	Config     *config.Config
	EtcdClient *etcd_client.EtcdClient
	Cache      *cache.Cache
}

/*type Callback func(rl *RuleLoader) error


var RollCallbackFunc  map[string]Callback

func init(){
	RollCallbackFunc = make(map[string]Callback)
}*/

func NewPromRuleLoader(cfg *config.Config) *RuleLoader {
	cli, err := etcd_client.NewClient(cfg.Etcd_addr)

	if err != nil {
		return nil
	}
	ca := cache.NewCache()
	prom_rule_loader := &RuleLoader{
		Config:     cfg,
		EtcdClient: cli,
		Cache:      ca,
	}
	return prom_rule_loader
}

func (rl *RuleLoader) Serve() error {
	prefix := "/prometheus/alert/rules/"
	rulelist, err := rl.EtcdClient.List(prefix)
	if err != nil {
		fmt.Println("etcdlist error...")
		return err
	}

	for _, key := range rulelist {
		value, err := rl.EtcdClient.Get(key)
		if err != nil {
			return err
		}

		status := rl.Cache.Set(key, value)

		if status != true {
			err = fmt.Errorf("cache set failed")
			return err
		}
	}

	http.HandleFunc("/append", func(w http.ResponseWriter, r *http.Request) {
		reqBuf, err := ioutil.ReadAll(r.Body)

		fmt.Println(string(reqBuf))
		if err != nil {
			fmt.Println(err)
			fmt.Println("1...")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = rl.Append(prefix, string(reqBuf))
		if err != nil {
			fmt.Println("2...")
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("OK"))
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		var rule Rule
		reqBuf, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal([]byte(reqBuf), &rule)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = rl.Delete(prefix, string(reqBuf))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("OK"))
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		var rule Rule
		reqBuf, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal([]byte(reqBuf), &rule)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = rl.Update(prefix, string(reqBuf))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("OK"))
	})
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		var reply Rule
		var replyrule []Rule
		replyrule = make([]Rule, 0)
		prefix := "/prometheus/alert/rules/"
		rulelist, err = rl.EtcdClient.List(prefix)
		if err != nil {
			fmt.Println("etcdlist error...")
			return
		}

		for _, key := range rulelist {
			value, err := rl.EtcdClient.Get(key)
			if err != nil {
				return
			}
			err = json.Unmarshal([]byte(value), &reply)

			if err != nil {
				return
			}
			replyrule = append(replyrule, reply)
		}
		data, _ := json.Marshal(replyrule)
		fmt.Fprintf(w, "%v", string(data))
	})

	addr := ":" + rl.Config.Server_port
	return http.ListenAndServe(addr, nil)
}

func (rl *RuleLoader) Append(prefix string, value string) error {
	var ruleall []Rule
	err := json.Unmarshal([]byte(value), &ruleall)
	if err != nil {
		fmt.Println("json.Unmarshal...")
		return err
	}
	for _, rule := range ruleall {
		key := prefix + rule.Namespace + rule.Server_name + rule.Alertname
		if !rl.Cache.Set(key, value) {
			fmt.Println("Cache Set failed")
			err := fmt.Errorf("%s Set failed\n", key)
			return err
		}
		err = rl.EtcdClient.Set(key, value)
		if err != nil {
			fmt.Println("etcd set failed..")
			return err
		}
		err = rl.OneMerge(&rule)
		if err != nil {
			return err
		}
		err = rl.Reload()
		if err != nil {
			return err
		}
	}
	return nil
}

func Link(str string) string {
	var str1 string
	strarr1 := strings.Split(str, "_")
	for _, value := range strarr1 {
		strarr2 := strings.Split(value, "-")
		for _, value1 := range strarr2 {
			str1 = str1 + value1
		}
	}
	return str1

}
func Encoding(rule *Rule) string {

	namespace := Link(rule.Namespace)
	if namespace == "" {
		return ""
	}

	server_name := Link(rule.Server_name)

	if server_name == "" {
		return ""
	}

	if rule.Alertname == "" {
		return ""
	}
	alert := "\nALERT " + namespace + server_name + rule.Alertname + "\n"

	if rule.Condition == "" {
		return ""
	}
	condition := "IF " + rule.Condition + "\n"

	if rule.Time == "" {
		return ""
	}
	time := "FOR " + rule.Time + "\n"

	label := "LABELS {\n"

	item_name := "  item_name = " + `"` + rule.Namespace + `|` + rule.Server_name + `|` + rule.Alertname + `",` + "\n"

	//fmt.Println("ITEM_NAME")
	//fmt.Println(item_name)
	trigername := "  triggername = " + `"` + rule.Trigername + `",` + "\n"

	value := "  value = " + `"` + "{{ $value }}" + `",` + "\n"

	if rule.Hostgroup == "" {
		return ""
	}

	hostgroup := "  hostgroup = " + `"` + rule.Hostgroup + `",` + "\n"

	templatename := "  templatename = " + `"` + rule.Templatename + `",` + "\n"

	if rule.Level == "" {
		return ""
	}
	level := "  level = " + `"` + rule.Level + `",` + "\n}\n"

	annotation := "ANNOTATIONS {\n"

	summary := "  summary = " + `"` + rule.Summary + `"` + ",\n"

	if rule.Description == "" {
		return ""
	}
	description := "  description = " + `"` + rule.Description + `"` + "\n}\n"

	str := alert + condition + time + label + value + item_name + trigername + hostgroup + templatename + level + annotation + summary + description

	return str
}

func (rl *RuleLoader) Delete(prefix string, value string) error {
	var rule Rule
	err := json.Unmarshal([]byte(value), &rule)
	if err != nil {
		fmt.Println("json.Unmarshal...")
		return err
	}
	key := prefix + rule.Namespace + rule.Server_name + rule.Alertname

	rl.Cache.Delete(key)
	err = rl.EtcdClient.Delete(key)
	if err != nil {
		fmt.Println("etcd delete faied...")
		return err
	}

	err = rl.AllMerge()
	if err != nil {
		return err
	}
	err = rl.Reload()
	if err != nil {
		return err
	}

	return nil
}

func (rl *RuleLoader) Update(prefix string, value string) error {
	var rule Rule
	err := json.Unmarshal([]byte(value), &rule)
	if err != nil {
		fmt.Println("json.Unmarshal...")
		return err
	}

	key := prefix + rule.Namespace + rule.Server_name + rule.Alertname

	if !rl.Cache.Update(key, value) {
		fmt.Println("Cache update failed...")
		err := fmt.Errorf("%s update failed\n", key)
		return err
	}

	err = rl.EtcdClient.Update(key, value)
	if err != nil {
		fmt.Println("etcd update failed...")
		return err
	}

	err = rl.AllMerge()
	if err != nil {
		return err
	}

	err = rl.Reload()
	if err != nil {
		return err
	}

	return nil
}

//add rule of
func (rl *RuleLoader) OneMerge(rule *Rule) error {
	str := Encoding(rule)

	if str == "" {
		err := fmt.Errorf("rule is null...")
		return err
	}
	buf, err := ioutil.ReadFile(rl.Config.RuleLocation)
	if err != nil {
		return err
	}

	newstr := string(buf) + string(str)
	err = ioutil.WriteFile(rl.Config.RuleLocation, []byte(newstr), 0x644)
	if err != nil {
		return err
	}
	return nil
}

//merge all of rules when delete/update rule
func (rl *RuleLoader) AllMerge() error {
	var buf string
	var err error
	var rule Rule
	if rl.Config.OldRuleLocation != "" {
		readbuf, err := ioutil.ReadFile(rl.Config.OldRuleLocation)
		if err != nil {
			return err
		}
		buf = string(readbuf)
	}
	for _, value := range rl.Cache.Rulemap {
		switch value.(type) {
		case string:
			err = json.Unmarshal([]byte(value.(string)), &rule)
			str := Encoding(&rule)
			buf = buf + str
			println(value)
		default:
			continue
		}
	}
	err = ioutil.WriteFile(rl.Config.RuleLocation, []byte(buf), 0x644)
	if err != nil {
		return err
	}
	return nil
}

//Reload for prometheus
func (rl *RuleLoader) Reload() error {
	client := &http.Client{}

	time.Sleep(300 * time.Millisecond)
	url := rl.Config.ReloadEndpoints + "/-/reload"
	req, err := http.NewRequest("POST", string(url), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	fmt.Printf("Prometheus Reload is %s\n", string(body))
	return nil
}
