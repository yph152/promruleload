package config

type Config struct {
	Server_port     string
	RuleLocation    string
	Etcd_addr       []string
	ReloadEndpoints string
	OldRuleLocation string
}
