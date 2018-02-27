package main

import (
	"flag"
	"fmt"
	"github.com/yph152/goproject/promeloadrules/config"
	"github.com/yph152/goproject/promeloadrules/ruleloader"
	"os"
	"strings"
)

var (
	serverport      = flag.String("server_port", "8890", "server_port of the prometheusRuleLoader reset endpoints")
	rulelocation    = flag.String("rulepath", "", "a rule_files: location in your prometheus config.")
	oldrulelocation = flag.String("oldrulepath", "", "a old_rule_files: location in your prometheus config for exist and unsimple rulelocation.")
	etcdaddr        = flag.String("etcdendpoints", "http://127.0.0.1:2379", "etcd_addr: for etcd server address.")
	reloadendpoints = flag.String("endpoint", "", "Endpoint of the prometheus reset endpoint (eg: http://prometheus:9090/-/reload).")
	helpFlag        = flag.Bool("help", false, "")
)

func main() {
	flag.Parse()

	if *helpFlag ||
		*rulelocation == "" ||
		*etcdaddr == "" ||
		*reloadendpoints == "" {
		fmt.Println("flag is nil...")
		os.Exit(0)
	}

	var etcdendpoints []string
	arr := strings.Split(*etcdaddr, ",")
	for _, value := range arr {
		etcdendpoints = append(etcdendpoints, value)
	}
	cfg := &config.Config{
		Server_port:     *serverport,
		RuleLocation:    *rulelocation,
		Etcd_addr:       etcdendpoints,
		ReloadEndpoints: *reloadendpoints,
		OldRuleLocation: *oldrulelocation,
	}

	fmt.Printf("Rule Update loaded.\n")
	fmt.Printf("Rule location: %s\n", *rulelocation)

	prom_rule_loader := ruleloader.NewPromRuleLoader(cfg)
	//	err := prom_rule_loader.RuleLoader()
	err := prom_rule_loader.Serve()

	if err != nil {
		fmt.Println(err)
		fmt.Println("Quit...")
		os.Exit(0)
	}
}
