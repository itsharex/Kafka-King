/*
 *
 *  * Copyright (c) 2025 Bronya0 <tangssst@163.com>. All rights reserved.
 *  * Original source: https://github.com/Bronya0
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *    http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package service

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"log"
	"testing"
)

// 引入 testing 包

func TestNewKafkaService2(t *testing.T) { // 功能测试以 `Test` 前缀命名

	ctx := context.Background()

	seeds := []string{"localhost:9092"}

	// 2. 自定义TLS配置
	// 加载CA证书
	//caCert, err := os.ReadFile("ca.pem")
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)
	//// 加载客户端证书和私钥
	//cert, err := tls.LoadX509KeyPair("client-cert.pem", "client-key.pem")
	//tlsConfig := &tls.Config{
	//	RootCAs:      caCertPool,
	//	Certificates: []tls.Certificate{cert},
	//	MinVersion:   tls.VersionTLS12,
	//}

	//创建客户端
	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.SASL(plain.Auth{
			User: "admin",
			Pass: "admin-secret",
		}.AsMechanism()),
		kgo.ProducerBatchCompression(kgo.GzipCompression()),
		kgo.ConsumerGroup("g1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	//包装admin客户端
	admin := kadm.NewClient(client)
	defer client.Close()

	//生产消息
	//st := time.Now()
	//for i := 0; i < 10; i++ {
	//	client.Produce(ctx, &kgo.Record{
	//		Topic: "1",
	//		Value: []byte(time.Now().Format(time.DateTime)),
	//	}, nil)
	//}
	//fmt.Printf("耗时：%.4f秒\n", time.Now().Sub(st).Seconds())

	////消费消息。订阅也可以在创建客户端的时候做
	client.AddConsumeTopics("1")
	fetches := client.PollRecords(context.Background(), 100)
	if errs := fetches.Errors(); len(errs) > 0 {
		panic(fmt.Sprint(errs))
	}
	for _, _ = range fetches.Records() {
		//fmt.Printf("%+v\n", string(record.Value))
	}

	////提交offset
	//if err := admin.CommitAllOffsets(ctx, "g1", kadm.OffsetsFromFetches(fetches)); err != nil {
	//	log.Fatalf("提交offsets失败: %v", err)
	//}
	//
	////读取offset
	//res2, err := admin.FetchOffsetsForTopics(ctx, "g1", "1")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//newMap := map[string]map[int32]any{}
	//for k1, v := range res2.Offsets() {
	//	if _, ok := newMap[k1]; !ok {
	//		newMap[k1] = make(map[int32]any)
	//	}
	//	for k2, v2 := range v {
	//		//fmt.Printf("%+v\n", v2)
	//		m := common.StructToMap(v2)
	//		newMap[k1][k2] = m
	//	}
	//}
	//fmt.Printf("%+v\n", newMap)

	//topics, err := admin.ListTopics(ctx)
	//fmt.Printf("%+v", topics)

	//fmt.Println(admin.ListBrokers(ctx))

	//fmt.Println(admin.ListGroups(ctx))

	//fmt.Println(admin.ListTopicsWithInternal(ctx))

	//res, _ := admin.ListStartOffsets(ctx, "1")
	//fmt.Printf("%+v\n", res)
	//
	//res, _ = admin.ListCommittedOffsets(ctx, "1")
	//fmt.Printf("%+v\n", res)
	//
	//res, _ = admin.ListEndOffsets(ctx, "1")
	//fmt.Printf("%+v\n", res)

	//ks.CreateTopics([]map[string]any{
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

	//map[g1:{Group:g1 Coordinator:{NodeID:0 Port:9092 Host:DESKTOP-7QTQFHC.mshome.net Rack:<nil> _:{}} State:Stable ProtocolType:consumer Protocol:cooperative-sticky
	//Members:[{MemberID:kgo-eb77103b-d127-4f0d-9159-6bdc92030cd1 InstanceID:<nil> ClientID:kgo ClientHost:/192.168.160.1 Join:{i:0xc000032960} Assigned:{i:0xc000284000}}] Err:<nil>}]
	res, _ := admin.DescribeGroups(ctx, "g1")
	fmt.Printf("%+v\n", res)
}
