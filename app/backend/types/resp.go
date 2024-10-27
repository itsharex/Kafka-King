package types

type Tag struct {
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
	Results []interface{} `json:"results"`
	Err     string        `json:"err"`
}
type ResultResp struct {
	Result map[string]interface{} `json:"result"`
	Err    string                 `json:"err"`
}
type Connect struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	BootstrapServers string `json:"bootstrap_servers"`
	Tls              string `json:"tls"`
	SkipTLSVerify    string `json:"skipTLSVerify"`
	TlsCertFile      string `json:"tls_cert_file"`
	TlsKeyFile       string `json:"tls_key_file"`
	TlsCaFile        string `json:"tls_ca_file"`
	Sasl             string `json:"sasl"`
	SaslMechanism    string `json:"sasl_mechanism"`
	SaslUser         string `json:"sasl_user"`
	SaslPwd          string `json:"sasl_pwd"`
}
type H map[string]interface{}

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
