package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"os"
	"testing"
)

// 引入 testing 包

func TestNewKafkaService2(t *testing.T) { // 功能测试以 `Test` 前缀命名
	//ks := NewKafkaService()
	//conn := map[string]interface{}{
	//	"bootstrap_servers": "127.0.0.1:9092",
	//	"tls":               "disable",
	//	"tls_skip_verify":   "false",
	//	"tls_ca_file":       "path/to/ca.pem",
	//	"tls_cert_file":     "path/to/client-cert.pem",
	//	"tls_key_file":      "path/to/client-key.pem",
	//	"sasl":              "enable",
	//	"sasl_mechanism":    "PLAIN",
	//	"sasl_user":         "admin",
	//	"sasl_pwd":          "admin-secret",
	//}
	//
	//err := ks.SetConnect("test", conn, false)
	//if err.Err != "" {
	//	panic(err.Err)
	//}
	ctx := context.Background()

	seeds := []string{"localhost:9092"}
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.

	//if err != nil {
	//	panic(err)
	//}
	//defer cl.Close()

	// 2. 自定义TLS配置
	// 加载CA证书
	caCert, err := os.ReadFile("ca.pem")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	// 加载客户端证书和私钥
	cert, err := tls.LoadX509KeyPair("client-cert.pem", "client-key.pem")
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	client, _ := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.SASL(plain.Auth{
			User: "admin",
			Pass: "admin-secret",
		}.AsMechanism()),
		kgo.DialTLSConfig(tlsConfig),
		//kgo.ConsumerGroup("my-group-identifier"),
		//kgo.ConsumeTopics("foo"),
	)

	admin := kadm.NewClient(client)
	topics, err := admin.ListTopics(ctx)
	if err != nil {
		return
	}
	// 计算每个主题的副本因子
	for _, topic := range topics {
		topicName := topic.Topic
		partitionCount := len(topic.Partitions)
		replicationFactor := len(topic.Partitions[0].Replicas)
	}
	fmt.Printf("%+v", topics)
	fmt.Println(admin.ListBrokers(ctx))
	fmt.Println(admin.ListGroups(ctx))
	fmt.Println(admin.ListTopicsWithInternal(ctx))

	//ks.CreateTopics([]map[string]interface{}{
	//	{
	//		"topic":             "test3",
	//		"numPartitions":     2,
	//		"replicationFactor": 1,
	//		"cleanupPolicy":     "delete",
	//		"retentionMs":       "604800000",
	//	},
	//})
	//res := ks.GetBrokers().Result

	//fmt.Printf("%+v", res["brokers"])
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
