package service

import (
	"app/backend/common"
	"app/backend/types"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/jcmturner/gokrb5/v8/client"
	krbConfig "github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/kerberos"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"log"
	"os"
	"strings"
	"sync"
	"time"
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
	mutex            sync.Mutex
}

func (k *Service) ptr(s string) *string {
	return &s
}

func NewKafkaService() *Service {
	return &Service{}
}

func (k *Service) SetConnect(connectName string, conn map[string]interface{}, isTest bool) *types.ResultResp {
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
			config = append(config, kgo.SASL(kerberos.Auth{
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

	config = append(config, kgo.SeedBrokers(bootstrapServers...))

	cl, err := kgo.NewClient(config...)
	if err != nil {
		result.Err = err.Error()
		return result
	}
	admin := kadm.NewClient(cl)
	ctx := context.Background()
	_, err = admin.ListTopics(ctx)
	if err != nil {
		log.Println("连接集群失败", err)
		result.Err = err.Error()
		return result
	}
	if isTest == false {
		k.connectName = connectName
		k.kac = admin
		k.client = cl
		k.config = config
		k.bootstrapServers = bootstrapServers
	}

	return result
}

// TestClient 测试连接
func (k *Service) TestClient(connectName string, conn map[string]interface{}) *types.ResultResp {
	return k.SetConnect(connectName, conn, true)
}

// GetBrokers 获取集群信息
func (k *Service) GetBrokers() *types.ResultResp {
	result := &types.ResultResp{}

	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()
	brokers, err := k.kac.ListBrokers(ctx)
	if err != nil {
		result.Err = err.Error()
		return result
	}

	var brokersResp []map[string]interface{}
	for _, broker := range brokers {
		brokersResp = append(brokersResp, map[string]interface{}{
			"node_id": broker.NodeID,
			"host":    broker.Host,
			"port":    broker.Port,
			"rack":    broker.Rack,
		})
	}

	clusterInfo := map[string]interface{}{
		"brokers": brokersResp,
	}
	result.Result = clusterInfo
	return result
}

// GetBrokerConfig 获取Broker配置
func (k *Service) GetBrokerConfig(brokerID int32) *types.ResultsResp {
	result := &types.ResultsResp{}

	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()

	configs, err := k.kac.DescribeBrokerConfigs(ctx, brokerID)
	if err != nil {
		result.Err = err.Error()
		return result
	}
	cfg := configs[0].Configs
	// 转换为map格式
	for _, config := range cfg {
		result.Results = append(result.Results, map[string]interface{}{
			"Name":      config.Key,
			"Value":     config.Value,
			"Source":    config.Source.String(),
			"Sensitive": config.Sensitive,
		})
	}
	return result
}

// GetTopics 获取主题信息
func (k *Service) GetTopics() *types.ResultsResp {
	result := &types.ResultsResp{}

	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()

	topics, err := k.kac.ListTopics(ctx)
	if err != nil {
		result.Err = err.Error()
		return result
	}

	for topicName, topicDetail := range topics {
		var partitions []interface{}
		for _, partition := range topicDetail.Partitions {
			errMsg := ""
			if partition.Err != nil {
				errMsg = partition.Err.Error()
			}
			partitions = append(partitions, map[string]interface{}{
				"partition":       partition.Partition,
				"leader":          partition.Leader,
				"replicas":        partition.Replicas,
				"isr":             partition.ISR,
				"err":             errMsg,
				"LeaderEpoch":     partition.LeaderEpoch,
				"OfflineReplicas": partition.OfflineReplicas,
			})
		}
		resultErrMsg := ""
		if topicDetail.Err != nil {
			resultErrMsg = topicDetail.Err.Error()
		}
		result.Results = append(result.Results, map[string]interface{}{
			"ID":                 topicDetail.ID,
			"topic":              topicName,
			"partition_count":    len(topicDetail.Partitions),
			"replication_factor": len(topicDetail.Partitions[0].Replicas),
			"IsInternal":         topicDetail.IsInternal,
			"Err":                resultErrMsg,
			"partitions":         partitions,
		})
	}

	return result
}

// GetTopicConfig 获取主题配置
func (k *Service) GetTopicConfig(topic string) *types.ResultsResp {
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()

	res, err := k.kac.DescribeTopicConfigs(ctx, topic)
	if err != nil {
		result.Err = err.Error()
		return result
	}
	cfg := res[0].Configs
	for _, config := range cfg {
		result.Results = append(result.Results, map[string]interface{}{
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
		result.Err = "请先选择连接"
		return result
	}

	ctx := context.Background()
	startOffsets, err := k.kac.ListStartOffsets(ctx, topics...)
	if err != nil {
		result.Err = err.Error()
		return result
	}

	endOffsets, err := k.kac.ListEndOffsets(ctx, topics...)
	if err != nil {
		result.Err = err.Error()
		return result
	}

	//读取offset
	committedOffsets, err := k.kac.FetchOffsetsForTopics(ctx, groupID, topics...)
	if err != nil {
		result.Err = err.Error()
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
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()
	groups, err := k.kac.ListGroups(ctx)
	if err != nil {
		result.Err = err.Error()
		return result
	}

	for group := range groups {
		result.Results = append(result.Results, map[string]interface{}{
			"Group":        group,
			"State":        groups[group].State,
			"ProtocolType": groups[group].ProtocolType,
			"Coordinator":  groups[group].Coordinator,
		})
	}

	return result
}

// GetGroupMembers 获取消费组下的成员信息
func (k *Service) GetGroupMembers(groupLst []string) *types.ResultsResp {
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()
	groups, err := k.kac.DescribeGroups(ctx, groupLst...)
	if err != nil {
		result.Err = err.Error()
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

//
//// DeleteConsumerGroup 删除消费组
//func (k *Service) DeleteConsumerGroup(groupID string) *types.ResultResp {
//	result := &types.ResultResp{}
//	if k.kac == nil {
//		result.Err = "请先选择连接"
//		return result
//	}
//	err := k.kac.DeleteConsumerGroup(groupID)
//	if err != nil {
//		result.Err = err.Error()
//		return result
//	}
//	return nil
//}

// CreateTopics 创建主题
func (k *Service) CreateTopics(topics []string, numPartitions, replicationFactor int, configs map[string]string) *types.ResultResp {
	result := &types.ResultResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	// 转换为 map[string]*string
	pointerMap := make(map[string]*string)
	for key, value := range configs {
		pointerMap[key] = &value
	}
	ctx := context.Background()
	_, err := k.kac.CreateTopics(ctx, int32(numPartitions), int16(replicationFactor), pointerMap, topics...)
	if err != nil {
		result.Err = err.Error()
		return result
	}

	return result
}

// DeleteTopic 删除主题
func (k *Service) DeleteTopic(topics []string) *types.ResultResp {
	result := &types.ResultResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()

	for _, topic := range topics {
		_, err := k.kac.DeleteTopic(ctx, topic)
		if err != nil {
			result.Err = err.Error()
			return result
		}
	}
	return result
}

// DeleteGroup 删除Group
func (k *Service) DeleteGroup(groups []string) *types.ResultResp {
	result := &types.ResultResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()
	_, err := k.kac.DeleteGroups(ctx, groups...)
	if err != nil {
		result.Err = err.Error()
		return result
	}
	return result
}

// CreatePartitions 添加分区
func (k *Service) CreatePartitions(topics []string, count int) *types.ResultResp {
	result := &types.ResultResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()

	for _, topic := range topics {
		_, err := k.kac.CreatePartitions(ctx, count, topic)
		if err != nil {
			result.Err = err.Error()
			return result
		}
	}

	return result
}

// AlterTopicConfig 修改主题配置
func (k *Service) AlterTopicConfig(topic string, name, value string) *types.ResultsResp {
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
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
	_, err := k.kac.AlterTopicConfigs(ctx, ac, topic)
	if err != nil {
		result.Err = err.Error()
		return result
	}
	return result
}

func (k *Service) AlterNodeConfig(nodeId int32, name, value string) *types.ResultsResp {
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
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
	_, err := k.kac.AlterBrokerConfigs(ctx, ac, nodeId)
	if err != nil {
		result.Err = err.Error()
		return result
	}
	return result
}

// Produce 生产消息
func (k *Service) Produce(topic string, key, value string, partition, num int, headers []map[string]string) *types.ResultsResp {
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}
	ctx := context.Background()
	st := time.Now()
	headers2 := make([]kgo.RecordHeader, len(headers))
	for i := 0; i < len(headers); i++ {
		headers2[i] = kgo.RecordHeader{
			Key:   headers[i]["key"],
			Value: []byte(headers[i]["value"]),
		}
	}
	for i := 0; i < num; i++ {
		k.client.Produce(ctx, &kgo.Record{
			Topic:     topic,
			Value:     []byte(value),
			Key:       []byte(key),
			Headers:   headers2,
			Partition: int32(partition),
		}, nil)
	}
	fmt.Printf("耗时：%.4f秒\n", time.Now().Sub(st).Seconds())

	return result
}

// Consumer 消费消息
func (k *Service) Consumer(topic string, group string, num, timeout int) *types.ResultsResp {
	result := &types.ResultsResp{}
	if k.kac == nil {
		result.Err = "请先选择连接"
		return result
	}

	st := time.Now()

	log.Println("消费消息。订阅也可以在创建客户端的时候做...")

	currentTopics := k.client.GetConsumeTopics()
	if len(currentTopics) == 1 && currentTopics[0] == topic {
		log.Println("当前消费主题和订阅主题一致，无需切换")
	} else {
		if len(currentTopics) > 0 {
			// 1. 清除所有当前正在消费的topics
			k.client.PurgeTopicsFromConsuming(currentTopics...)
		}
		// 2. 添加新的topics
		k.client.AddConsumeTopics(topic)
		//k.client.ResumeFetchTopics(topic)
	}

	log.Println("开始poll...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	fetches := k.client.PollRecords(ctx, num)
	cancel() // 注意要调用cancel避免context泄露

	if errs := fetches.Errors(); len(errs) > 0 {
		result.Err = fmt.Sprint(errs)
		return result
	}
	res := make([]any, 0)
	log.Println("poll完成...", len(fetches.Records()))
	for i, v := range fetches.Records() {
		if v == nil {
			continue
		}
		res = append(res, map[string]any{
			"ID":            i,
			"Offset":        v.Offset,
			"Key":           string(v.Key),
			"Value":         string(v.Value),
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

	fmt.Printf("耗时：%.4f秒\n", time.Now().Sub(st).Seconds())
	fmt.Println(topic, group, num)

	log.Println("提交offset...")
	if group != "" {
		if err := k.kac.CommitAllOffsets(context.Background(), group, kadm.OffsetsFromFetches(fetches)); err != nil {
			result.Err = "提交offsets失败: " + err.Error()
			return result
		}
	}
	return result
}
