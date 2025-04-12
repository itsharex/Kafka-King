/*
 * Copyright 2025 Bronya0 <tangssst@163.com>.
 * Author Github: https://github.com/Bronya0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"app/backend/common"
	"app/backend/types"
	"app/backend/utils/compress"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jcmturner/gokrb5/v8/client"
	krbConfig "github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/kerberos"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

type TopicConfig struct {
	Name              string
	NumPartitions     int32
	ReplicationFactor int16
}

type Service struct {
	connectName      string
	bootstrapServers []string
	config           []kgo.Opt
	kac              *kadm.Client
	client           *kgo.Client
	consumer         []any
	mutex            sync.Mutex
	topics           []any
	groups           []any
}

func (k *Service) ptr(s string) *string {
	return &s
}

func NewKafkaService() *Service {
	return &Service{}
}

func (k *Service) Close(ctx context.Context) bool {
	// 关闭连接,return false表示允许关闭
	if k.client != nil {
		k.client.Close()
	}
	if k.kac != nil {
		k.kac.Close()
	}
	if k.consumer != nil && len(k.consumer) == 2 {
		k.consumer[2].(*kgo.Client).Close()
	}
	fmt.Println("关闭连接")
	return false
}

func (k *Service) SetConnect(connectName string, conn map[string]any, isTest bool) *types.ResultResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultResp{}

	var config []kgo.Opt

	// TLS配置
	if conn["tls"] == "enable" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: conn["skipTLSVerify"] == "true", // 开发环境可以设置为true
		}

		// 如果需要证书认证
		if conn["tls_cert_file"] != "" && conn["tls_key_file"] != "" {
			cert, err := tls.LoadX509KeyPair(conn["tls_cert_file"].(string), conn["tls_key_file"].(string))
			if err != nil {
				log.Println("loading x509 key pair failed:", err)
				result.Err = fmt.Sprintf("loading x509 key pair failed: %v", err)
				return result

			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// 如果需要CA证书
		if conn["tls_ca_file"] != "" {
			caCert, err := os.ReadFile(conn["tls_ca_file"].(string))
			if err != nil {
				log.Println("reading CA file failed:", err)
				result.Err = fmt.Sprintf("reading CA file failed: %v", err)
				return result
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}

		config = append(config, kgo.DialTLSConfig(tlsConfig))
	}

	// SASL配置
	if conn["sasl"] == "enable" {
		user := conn["sasl_user"].(string)
		pwd := conn["sasl_pwd"].(string)

		// SASL机制设置
		mechanism := conn["sasl_mechanism"].(string)
		switch strings.ToUpper(mechanism) {
		case "PLAIN":
			config = append(config, kgo.SASL(plain.Auth{User: user, Pass: pwd}.AsMechanism()))
		case "SCRAM-SHA-256":
			config = append(config, kgo.SASL(scram.Auth{User: user, Pass: pwd}.AsSha256Mechanism()))
		case "SCRAM-SHA-512":
			config = append(config, kgo.SASL(scram.Auth{User: user, Pass: pwd}.AsSha512Mechanism()))
		case "GSSAPI":
			// 创建Kerberos配置
			// 1. 首先读取keytab文件
			kt, err := keytab.Load(conn["kerberos_user_keytab"].(string))
			if err != nil {
				result.Err = err.Error()
				return result
			}
			// 2. 读取krb5.conf配置
			cfg, err := krbConfig.Load(conn["kerberos_krb5_conf"].(string))
			if err != nil {
				result.Err = err.Error()
				return result
			}
			// 3. 创建客户端
			kerberosClient := client.NewWithKeytab(
				conn["Kerberos_user"].(string),  // username (principal的第一部分)
				conn["Kerberos_realm"].(string), // realm (Kerberos领域，大写的域名)
				kt,                              // keytab对象
				cfg,                             // krb5配置对象
				client.DisablePAFXFAST(true),    // 禁用PA-FX-FAST，提高兼容性
			)
			// 创建GSSAPI认证
			config = append(config,
				kgo.SASL(kerberos.Auth{
					Client:           kerberosClient,
					Service:          conn["kerberos_service_name"].(string),
					PersistAfterAuth: true, // 保留上下文
				}.AsMechanism()))

		default:
			log.Println("不支持的SASL机制", mechanism)
			result.Err = fmt.Sprintf("unsupported SASL mechanism: %s", mechanism)
			return result
		}
	}
	bootstrapServers := strings.Split(conn["bootstrap_servers"].(string), ",")

	config = append(
		config,
		kgo.SeedBrokers(bootstrapServers...),
	)

	cl, err := kgo.NewClient(config...)
	if err != nil {
		result.Err = "NewClient Error：" + err.Error()
		return result
	}
	admin := kadm.NewClient(cl)
	ctx := context.Background()
	topics, err := admin.ListTopics(ctx)
	if err != nil {
		log.Println("连接集群失败", err)
		result.Err = "ListTopics Error：" + err.Error()
		return result
	}

	//正式切换节点，赋值并清理缓存，并更新缓存
	if !isTest {
		k.connectName = connectName
		k.kac = admin
		k.client = cl
		k.config = config
		k.consumer = nil
		k.bootstrapServers = bootstrapServers
		k.clearCache()
		k.topics = k.buildTopicsResp(topics)
	}

	return result
}

// TestClient 测试连接
func (k *Service) TestClient(connectName string, conn map[string]any) *types.ResultResp {
	return k.SetConnect(connectName, conn, true)
}

func (k *Service) clearCache() {
	k.topics = nil
	k.groups = nil
}

// GetBrokers 获取集群信息
func (k *Service) GetBrokers() *types.ResultResp {
	result := &types.ResultResp{}

	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}
	ctx := context.Background()
	brokers, err := k.kac.ListBrokers(ctx)
	if err != nil {
		result.Err = "ListBrokers Error：" + err.Error()
		return result
	}

	var brokersResp []map[string]any
	for _, broker := range brokers {
		brokersResp = append(brokersResp, map[string]any{
			"node_id": broker.NodeID,
			"host":    broker.Host,
			"port":    broker.Port,
			"rack":    broker.Rack,
		})
	}

	clusterInfo := map[string]any{
		"brokers": brokersResp,
	}
	result.Result = clusterInfo
	return result
}

// GetBrokerConfig 获取Broker配置
func (k *Service) GetBrokerConfig(brokerID int32) *types.ResultsResp {
	result := &types.ResultsResp{}

	if k.kac == nil {
		result.Err = common.PleaseSelectErr // 确保 common.PleaseSelectErr 是一个有效的错误变量
		return result
	}
	ctx := context.Background()

	configs, err := k.kac.DescribeBrokerConfigs(ctx, brokerID)
	if err != nil {
		result.Err = fmt.Sprintf("DescribeBrokerConfigs Error：%s", err.Error())
		return result
	}

	if len(configs) == 0 {
		result.Err = "No configurations found for the given broker ID"
		return result
	}

	cfg := configs[0].Configs
	// 转换为map格式
	for _, config := range cfg {
		result.Results = append(result.Results, map[string]any{
			"Name":      config.Key,
			"Value":     config.Value,
			"Source":    config.Source.String(),
			"Sensitive": config.Sensitive,
		})
	}
	return result
}

func (k *Service) buildTopicsResp(topics kadm.TopicDetails) []any {
	var result []any
	for topicName, topicDetail := range topics {
		partitionErrs := ""
		var partitions []any
		for _, partition := range topicDetail.Partitions {
			errMsg := ""
			if partition.Err != nil {
				errMsg = partition.Err.Error()
				partitionErrs += fmt.Sprintf("partition %d: %s\n", partition.Partition, errMsg)
			}
			partitions = append(partitions, map[string]any{
				"partition":       partition.Partition,
				"leader":          partition.Leader,
				"replicas":        partition.Replicas,
				"isr":             partition.ISR,
				"err":             errMsg,
				"LeaderEpoch":     partition.LeaderEpoch,
				"OfflineReplicas": partition.OfflineReplicas,
			})
		}
		if topicDetail.Err != nil {
			partitionErrs = topicDetail.Err.Error() + "\n" + partitionErrs
		}
		// 检查分区列表是否为空，避免访问空切片的第一个元素
		replicationFactor := 0
		if len(topicDetail.Partitions) > 0 {
			replicationFactor = len(topicDetail.Partitions[0].Replicas)
		}
		result = append(result, map[string]any{
			"ID":                 topicDetail.ID,
			"topic":              topicName,
			"partition_count":    len(topicDetail.Partitions),
			"replication_factor": replicationFactor,
			"IsInternal":         topicDetail.IsInternal,
			"Err":                partitionErrs,
			"partitions":         partitions,
		})
	}
	return result
}

// GetTopics 获取主题信息
func (k *Service) GetTopics(noCache bool) *types.ResultsResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultsResp{}

	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}

	if !noCache && len(k.topics) > 0 {
		result.Results = k.topics
		return result
	}

	ctx := context.Background()

	topics, err := k.kac.ListTopics(ctx)
	if err != nil {
		result.Err = fmt.Sprintf("ListTopics Error：%v", err.Error())
		return result
	}
	result.Results = k.buildTopicsResp(topics)

	k.topics = result.Results
	return result
}

// GetTopicConfig 获取主题配置
func (k *Service) GetTopicConfig(topic string) *types.ResultsResp {
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}
	ctx := context.Background()

	res, err := k.kac.DescribeTopicConfigs(ctx, topic)
	if err != nil {
		result.Err = err.Error()
		return result
	}
	if len(res) == 0 {
		result.Err = "No configurations found for the given topic"
		return result
	}
	cfg := res[0].Configs
	for _, config := range cfg {
		result.Results = append(result.Results, map[string]any{
			"Name":      config.Key,
			"Value":     config.Value,
			"Source":    config.Source,
			"Synonyms":  config.Synonyms,
			"Sensitive": config.Sensitive,
		})
	}
	return result
}

func (k *Service) GetTopicOffsets(topics []string, groupID string) *types.ResultResp {
	result := &types.ResultResp{}

	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}

	ctx := context.Background()
	startOffsets, err := k.kac.ListStartOffsets(ctx, topics...)
	if err != nil {
		result.Err = "ListStartOffsets Error：" + err.Error()
		return result
	}

	endOffsets, err := k.kac.ListEndOffsets(ctx, topics...)
	if err != nil {
		result.Err = "ListEndOffsets Error：" + err.Error()
		return result
	}

	//读取offset
	committedOffsets, err := k.kac.FetchOffsetsForTopics(ctx, groupID, topics...)
	if err != nil {
		result.Err = "FetchOffsetsForTopics Error：" + err.Error()
		return result
	}

	// {"topicname":{"0":{"Topic":"1","Partition":0,"At":100,"LeaderEpoch":0,"Metadata":""},"1":。。。
	result.Result = map[string]any{
		"start_map":  k.ToMap(startOffsets.Offsets()),
		"end_map":    k.ToMap(endOffsets.Offsets()),
		"commit_map": k.ToMap(committedOffsets.Offsets()),
	}

	return result
}
func (k *Service) ToMap(mapStruct map[string]map[int32]kadm.Offset) map[string]map[int32]any {
	newMap := map[string]map[int32]any{}
	for k1, v := range mapStruct {
		if _, ok := newMap[k1]; !ok {
			newMap[k1] = make(map[int32]any)
		}
		for k2, v2 := range v {
			m := common.StructToMap(v2)
			newMap[k1][k2] = m
		}
	}
	return newMap
}

// GetGroups 获取消费组信息
func (k *Service) GetGroups() *types.ResultsResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}
	if len(k.groups) > 0 {
		result.Results = k.groups
		return result
	}
	ctx := context.Background()
	groups, err := k.kac.ListGroups(ctx)
	if err != nil {
		result.Err = "ListGroups Error：" + err.Error()
		return result
	}

	for group := range groups {
		result.Results = append(result.Results, map[string]any{
			"Group":        group,
			"State":        groups[group].State,
			"ProtocolType": groups[group].ProtocolType,
			"Coordinator":  groups[group].Coordinator,
		})
	}
	k.groups = result.Results
	return result
}

// GetGroupMembers 获取消费组下的成员信息
func (k *Service) GetGroupMembers(groupLst []string) *types.ResultsResp {
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}
	ctx := context.Background()
	groups, err := k.kac.DescribeGroups(ctx, groupLst...)
	if err != nil {
		result.Err = "DescribeGroups Error：" + err.Error()
		return result
	}
	//map[g1:{Group:g1 Coordinator:{NodeID:0 Port:9092 Host:DESKTOP-7QTQFHC.mshome.net Rack:<nil> _:{}} State:Stable ProtocolType:consumer Protocol:cooperative-sticky
	//Members:[{MemberID:kgo-eb77103b-d127-4f0d-9159-6bdc92030cd1 InstanceID:<nil> ClientID:kgo ClientHost:/192.168.160.1 Join:{i:0xc000032960} Assigned:{i:0xc000284000}}] Err:<nil>}]
	for key := range groups {
		members := groups[key].Members
		membersLst := make([]any, 0)
		for _, member := range members {
			membersLst = append(membersLst, map[string]any{
				"MemberID":   member.MemberID,
				"InstanceID": member.InstanceID,
				"ClientID":   member.ClientID,
				"ClientHost": member.ClientHost,
			})

		}

		result.Results = append(result.Results, map[string]any{
			"Group":        key,
			"Coordinator":  groups[key].Coordinator.Host,
			"State":        groups[key].State,
			"ProtocolType": groups[key].ProtocolType,
			"Protocol":     groups[key].Protocol,
			"Members":      membersLst,
		})
	}

	return result
}

// CreateTopics 创建主题
func (k *Service) CreateTopics(topics []string, numPartitions, replicationFactor int, configs map[string]string) *types.ResultResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}

	k.clearCache()

	// 转换为 map[string]*string
	pointerMap := make(map[string]*string)
	for key, value := range configs {
		pointerMap[key] = &value
	}

	ctx := context.Background()
	resp, err := k.kac.CreateTopics(ctx, int32(numPartitions), int16(replicationFactor), pointerMap, topics...)
	if err != nil {
		result.Err = "CreateTopics Error：" + err.Error()
		return result
	}
	if resp.Error() != nil {
		result.Err = "CreateTopics Error：" + resp.Error().Error()
		return result
	}

	return result
}

// DeleteTopic 删除主题
func (k *Service) DeleteTopic(topics []string) *types.ResultResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}

	k.clearCache()

	ctx := context.Background()
	for _, topic := range topics {
		resp, err := k.kac.DeleteTopic(ctx, topic)
		if err != nil {
			result.Err = "DeleteTopic Error：" + err.Error()
			return result
		}
		if resp.Err != nil {
			result.Err = "DeleteTopic Error：" + resp.Err.Error()
			return result
		}
	}
	return result
}

// DeleteGroup 删除Group
func (k *Service) DeleteGroup(group string) *types.ResultResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}

	k.clearCache()

	ctx := context.Background()
	resp, err := k.kac.DeleteGroup(ctx, group)
	if err != nil {
		result.Err = "DeleteGroup Error：" + err.Error()
		return result
	}
	if resp.Err != nil {
		result.Err = "DeleteGroup Error：" + resp.Err.Error()
		return result
	}
	return result
}

// CreatePartitions 添加分区
func (k *Service) CreatePartitions(topics []string, count int) *types.ResultResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}

	k.clearCache()

	ctx := context.Background()
	for _, topic := range topics {
		resp, err := k.kac.CreatePartitions(ctx, count, topic)
		if err != nil {
			result.Err = "CreatePartitions Error：" + err.Error()
			return result
		}
		if resp.Error() != nil {
			result.Err = "CreatePartitions Error：" + resp.Error().Error()
			return result
		}
	}

	return result
}

// AlterTopicConfig 修改主题配置
func (k *Service) AlterTopicConfig(topic string, name, value string) *types.ResultsResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}
	ac := []kadm.AlterConfig{
		{
			Name:  name,
			Op:    kadm.SetConfig,
			Value: &value,
		},
	}

	ctx := context.Background()
	resp, err := k.kac.AlterTopicConfigs(ctx, ac, topic)
	if err != nil {
		result.Err = "AlterTopicConfigs Error：" + err.Error()
		return result
	}
	for _, v := range resp {
		if v.Err != nil {
			result.Err = "AlterTopicConfigs Error：" + v.Err.Error()
			return result
		}
	}
	return result
}

func (k *Service) AlterNodeConfig(nodeId int32, name, value string) *types.ResultsResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}
	ac := []kadm.AlterConfig{
		{
			Name:  name,
			Op:    kadm.SetConfig,
			Value: &value,
		},
	}

	ctx := context.Background()
	resp, err := k.kac.AlterBrokerConfigs(ctx, ac, nodeId)
	if err != nil {
		result.Err = "AlterBrokerConfigs Error：" + err.Error()
		return result
	}
	for _, v := range resp {
		if v.Err != nil {
			result.Err = "AlterBrokerConfigs Error：" + v.Err.Error()
			return result
		}
	}
	return result
}

// Produce 生产消息
func (k *Service) Produce(topic string, key, value string, partition, num int, headers []map[string]string, compressMethod string) *types.ResultsResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	st := time.Now()
	headers2 := make([]kgo.RecordHeader, len(headers))
	for i := 0; i < len(headers); i++ {
		headers2[i] = kgo.RecordHeader{
			Key:   headers[i]["key"],
			Value: []byte(headers[i]["value"]),
		}
	}

	var data []byte
	var err error
	switch compressMethod {
	case "gzip":
		data, err = compress.Gzip([]byte(value))
	case "lz4":
		data, err = compress.Lz4([]byte(value))
	case "zstd":
		data, err = compress.Zstd([]byte(value))
	case "snappy":
		data, err = compress.Snappy([]byte(value))
	default:
		data = []byte(value)
	}
	if err != nil {
		result.Err = "Failed to compress data: " + err.Error()
		return result
	}
	var records []*kgo.Record
	for i := 0; i < num; i++ {
		records = append(records, &kgo.Record{
			Topic:     topic,
			Value:     data,
			Key:       []byte(key),
			Headers:   headers2,
			Partition: int32(partition),
		})
	}

	//同步发送
	res := k.client.ProduceSync(ctx, records...)
	if err := res.FirstErr(); err != nil {
		result.Err = "Produce Error：" + err.Error()
		return result
	}

	fmt.Printf("耗时：%.4f秒\n", time.Since(st).Seconds())

	return result
}

// Consumer 消费消息
// 参考：https://github.com/twmb/franz-go/blob/master/examples/group_consuming/main.go
func (k *Service) Consumer(topic string, group string, num, timeout int, decompress string, isCommit bool) *types.ResultsResp {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = common.PleaseSelectErr
		return result
	}
	st := time.Now()

	var _client *kgo.Client
	if k.consumer == nil || (k.consumer != nil && (k.consumer[0] != group || k.consumer[1] != topic)) {
		// 看看缓存的有没有，没有则关闭之前的，新缓存；有则用
		fmt.Println("创建新的consumer", k.consumer)
		if k.consumer != nil && len(k.consumer) == 3 {
			k.consumer[2].(*kgo.Client).Close()
		}
		conf := append(k.config,
			kgo.ConsumeTopics(topic),
			kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		)
		if group != "" {
			conf = append(conf, kgo.ConsumerGroup(group), kgo.DisableAutoCommit())
		}

		cl, err := kgo.NewClient(
			conf...,
		)
		if err != nil {
			result.Err = "Consumer Error：" + err.Error()
			return result
		}
		_client = cl
		k.consumer = []any{group, topic, _client}
		// 要等待重平衡成功。
		fmt.Println("创建新的consumer完成", k.consumer)
	} else {
		fmt.Println("使用缓存的consumer", k.consumer)
		_client = k.consumer[2].(*kgo.Client)
	}

	log.Println("开始poll...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	fetches := _client.PollRecords(ctx, num)

	// 客户端此时已经被关闭
	if fetches.IsClientClosed() {
		k.consumer = nil
		result.Err = "Client Closed, Please Retry"
		return result
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		result.Err = "Consume Timeout, Maybe No Message"
		return result
	}
	if errs := fetches.Errors(); len(errs) > 0 {
		result.Err = fmt.Sprint(errs)
		return result
	}
	log.Println("poll完成...", len(fetches.Records()))

	res := make([]any, 0)
	for i, v := range fetches.Records() {
		if v == nil {
			continue
		}

		var data []byte
		var err error
		switch decompress {
		case "gzip":
			data, err = compress.GzipDecompress(v.Value)
		case "lz4":
			data, err = compress.Lz4Decompress(v.Value)
		case "zstd":
			data, err = compress.ZstdDecompress(v.Value)
		case "snappy":
			data, err = compress.SnappyDecompress(v.Value)
		default:
			data = v.Value
		}
		if err != nil {
			result.Err = "Failed to decompress data: " + err.Error()
			return result
		}

		res = append(res, map[string]any{
			"ID":            i,
			"Offset":        v.Offset,
			"Key":           string(v.Key),
			"Value":         string(data),
			"Timestamp":     v.Timestamp.Format(time.DateTime),
			"Partition":     v.Partition,
			"Topic":         v.Topic,
			"Headers":       v.Headers,
			"LeaderEpoch":   v.LeaderEpoch,
			"ProducerEpoch": v.ProducerEpoch,
			"ProducerID":    v.ProducerID,
		})
	}
	result.Results = res

	fmt.Printf("耗时：%.4f秒\n", time.Since(st).Seconds())
	fmt.Println(topic, group, num)

	if group != "" && isCommit {
		log.Println("提交offset...")
		if err := _client.CommitUncommittedOffsets(context.Background()); err != nil {
			result.Err = "Failed to submit offsets: " + err.Error()
			return result
		}
	}
	return result

}
