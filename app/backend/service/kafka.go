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
)

const KingGroup = "kafka-king-group"

type TopicConfig struct {
	Name              string
	NumPartitions     int32
	ReplicationFactor int16
}

type Service struct {
	connectName      string
	bootstrapServers []string
	//config           *sarama.Config
	//kac              sarama.ClusterAdmin
	//consumer         sarama.Consumer
	//mutex            sync.Mutex
	config []kgo.Opt
	kac    *kadm.Client
	//consumer         sarama.Consumer
	mutex sync.Mutex
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
			)
			// 创建GSSAPI认证
			config = append(config, kgo.SASL(kerberos.Auth{
				Client:  kerberosClient,
				Service: conn["kerberos_service_name"].(string),
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
		k.config = config
		k.bootstrapServers = bootstrapServers
	}

	// Convert conn map to proper config
	// Add necessary configurations from conn map to sarama config
	//bootstrapServers := strings.Split(conn["bootstrap_servers"].(string), ",")
	//admin, err := sarama.NewClusterAdmin(bootstrapServers, config)
	//if err != nil {
	//	log.Println("创建Admin失败", err)
	//	result.Err = err.Error()
	//	return result
	//} else {
	//	if isTest == false {
	//		k.connectName = connectName
	//		k.kac = admin
	//		k.config = config
	//		k.bootstrapServers = bootstrapServers
	//	} else {
	//		_, err = admin.ListTopics()
	//		if err != nil {
	//			log.Println("连接集群失败", err)
	//			result.Err = err.Error()
	//			return result
	//		}
	//	}
	//}

	return result
}

// 创建Consumer
//func (k *Service) newConsumer() (sarama.Consumer, error) {
//	consumer, err := sarama.NewConsumer(k.bootstrapServers, k.config)
//	if err != nil {
//		return nil, err
//	}
//	return consumer, nil
//}

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

// // DescribeGroup 消费组详情
//
//	func (k *Service) DescribeGroup(groupID string) (*types.GroupInfo, error) {
//		// 获取消费组详情
//		groups, err := k.kac.DescribeConsumerGroups([]string{groupID})
//		if err != nil {
//			return nil, fmt.Errorf("describe consumer group failed: %v", err)
//		}
//		if len(groups) == 0 {
//			return nil, fmt.Errorf("group not found: %s", groupID)
//		}
//
//		// 获取消费组offset信息
//		offsetFetch, err := k.kac.ListConsumerGroupOffsets(groupID, nil)
//		if err != nil {
//			return nil, fmt.Errorf("list consumer group offsets failed: %v", err)
//		}
//
//		info := &types.GroupInfo{
//			Group:  groupID,
//			Topics: make(map[string][]types.PartitionOffset),
//		}
//
//		// 遍历每个topic的分区offset
//		for topic, partitions := range offsetFetch.Blocks {
//			var partitionOffsets []types.PartitionOffset
//
//			// 获取topic的最新offset
//			latestOffsets, err := k.getTopicLatestOffsets(topic)
//			if err != nil {
//				return nil, err
//			}
//
//			for partition, offsetBlock := range partitions {
//				if offsetBlock.Offset == -1 {
//					continue // 跳过未消费的分区
//				}
//
//				latestOffset := latestOffsets[partition]
//				lag := latestOffset - offsetBlock.Offset
//
//				po := types.PartitionOffset{
//					Partition: partition,
//					Offset:    offsetBlock.Offset,
//					Lag:       lag,
//				}
//				partitionOffsets = append(partitionOffsets, po)
//				info.TotalLag += lag
//			}
//			info.Topics[topic] = partitionOffsets
//		}
//
//		return info, nil
//	}

//func (k *Service) getLatestOffset(topic string, partition int32) (int64, error) {
//
//	client, err := sarama.NewClient(k.bootstrapServers, k.config)
//	if err != nil {
//		return 0, err
//	}
//	defer client.Close()
//
//	offset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
//	if err != nil {
//		return 0, err
//	}
//	return offset, nil
//}
