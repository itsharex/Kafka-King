package service

import (
	"fmt"
	"testing"
)

// 引入 testing 包

type a struct {
	name string
}

func TestNewKafkaService(t *testing.T) { // 功能测试以 `Test` 前缀命名
	ks := NewKafkaService()
	conn := map[string]interface{}{
		"bootstrap_servers": "127.0.0.1:9092",
		"tls":               "disable",
		"tls_skip_verify":   "false",
		"tls_ca_file":       "path/to/ca.pem",
		"tls_cert_file":     "path/to/client-cert.pem",
		"tls_key_file":      "path/to/client-key.pem",
		"sasl":              "enable",
		"sasl_mechanism":    "PLAIN",
		"sasl_user":         "admin",
		"sasl_pwd":          "admin-secret",
	}

	err := ks.SetConnect("test", conn, false)
	if err.Err != "" {
		panic(err.Err)
	}
	//ks.CreateTopics([]map[string]interface{}{
	//	{
	//		"topic":             "test3",
	//		"numPartitions":     2,
	//		"replicationFactor": 1,
	//		"cleanupPolicy":     "delete",
	//		"retentionMs":       "604800000",
	//	},
	//})
	res := ks.GetBrokers().Result

	fmt.Printf("%+v", res["brokers"])
	//
	//fmt.Println(ks.CreatePartitions([]string{"test3"}, 1))
	//fmt.Println(ks.GetTopics())
	//fmt.Println(ks.GetGroups())
	//fmt.Println(ks.GetTopicConfig("test3"))
	//fmt.Println(ks.DeleteTopic([]string{"test1"}))
	//fmt.Println(ks.DescribeTopic([]string{"test3"}))
	//fmt.Println(ks.AlterTopicConfig("test3", map[string]*string{"retention.ms": ptr("60480")}))
	//fmt.Println(ks.GetTopicConfig("test3"))
	//fmt.Println(ks.DescribeGroup("test"))
}
func ptr(s string) *string {
	return &s
}
