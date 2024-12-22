package types

type Tag struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}
type Config struct {
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Language string    `json:"language"`
	Theme    string    `json:"theme"`
	Connects []Connect `json:"connects"`
}
type ResultsResp struct {
	Results []any  `json:"results"`
	Err     string `json:"err"`
}
type ResultResp struct {
	Result map[string]any `json:"result"`
	Err    string         `json:"err"`
}
type Connect struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	BootstrapServers    string `json:"bootstrap_servers"`
	Tls                 string `json:"tls"`
	SkipTLSVerify       string `json:"skipTLSVerify"`
	TlsCertFile         string `json:"tls_cert_file"`
	TlsKeyFile          string `json:"tls_key_file"`
	TlsCaFile           string `json:"tls_ca_file"`
	Sasl                string `json:"sasl"`
	SaslMechanism       string `json:"sasl_mechanism"`
	SaslUser            string `json:"sasl_user"`
	SaslPwd             string `json:"sasl_pwd"`
	KerberosUserKeytab  string `json:"kerberos_user_keytab"`
	KerberosKrb5Conf    string `json:"kerberos_krb5_conf"`
	KerberosUser        string `json:"Kerberos_user"`
	KerberosRealm       string `json:"Kerberos_realm"`
	KerberosServiceName string `json:"kerberos_service_name"`
}
type H map[string]any

type GroupInfo struct {
	Group    string
	Topics   map[string][]PartitionOffset
	TotalLag int64
}

type PartitionOffset struct {
	Partition int32
	Offset    int64
	Lag       int64
}

// BrokerConfig Broker配置结构
type BrokerConfig struct {
	Name      string          `json:"name"`
	Value     string          `json:"value"`
	Source    string          `json:"source"` // DYNAMIC_BROKER_CONFIG, DEFAULT_CONFIG, STATIC_BROKER_CONFIG
	ReadOnly  bool            `json:"read_only"`
	Sensitive bool            `json:"sensitive"`
	Synonyms  []ConfigSynonym `json:"synonyms,omitempty"`
}

type ConfigSynonym struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Source string `json:"source"`
}
