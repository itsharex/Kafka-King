package service

import (
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
			"Name":  config.Key,
			"Value": config.Value,
			//"ReadOnly":  config.,
			//"Default":   config.,
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

	topics, err := k.kac.ListTopicsWithInternal(ctx)
	if err != nil {
		result.Err = err.Error()
		return result
	}

	for topicName, topicDetail := range topics {
		var partitions []interface{}
		for _, partition := range topicDetail.Partitions {
			partitions = append(partitions, map[string]interface{}{
				"partition":       partition.Partition,
				"leader":          partition.Leader,
				"replicas":        partition.Replicas,
				"isr":             partition.ISR,
				"err":             partition.Err.Error(),
				"LeaderEpoch":     partition.LeaderEpoch,
				"OfflineReplicas": partition.OfflineReplicas,
			})
		}
		result.Results = append(result.Results, map[string]interface{}{
			"ID":                 topicDetail.ID,
			"topic":              topicName,
			"partition_count":    len(topicDetail.Partitions),
			"replication_factor": len(topicDetail.Partitions[0].Replicas),
			"IsInternal":         topicDetail.IsInternal,
			"Err":                topicDetail.Err,
			"partitions":         partitions,
		})
	}

	return result
}

//
//// GetGroups 获取消费组信息
//func (k *Service) GetGroups() *types.ResultsResp {
//	result := &types.ResultsResp{}
//	if k.kac == nil {
//		result.Err = "请先选择连接"
//		return result
//	}
//
//	groups, err := k.kac.ListConsumerGroups()
//	if err != nil {
//		result.Err = err.Error()
//		return result
//	}
//
//	for groupName := range groups {
//		result.Results = append(result.Results, groupName)
//	}
//
//	return result
//}
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
//
//// CreateTopics 创建主题
//func (k *Service) CreateTopics(configs []map[string]interface{}) *types.ResultResp {
//	result := &types.ResultResp{}
//	if k.kac == nil {
//		result.Err = "请先选择连接"
//		return result
//	}
//
//	for _, config := range configs {
//		topic := config["topic"].(string)
//		topicDetail := &sarama.TopicDetail{
//			NumPartitions:     int32(config["numPartitions"].(int)),
//			ReplicationFactor: int16(config["replicationFactor"].(int)),
//			ConfigEntries: map[string]*string{
//				"cleanup.policy": k.ptr(config["cleanupPolicy"].(string)), // 或 "compact"
//				"retention.ms":   k.ptr(config["retentionMs"].(string)),   // 7天
//			},
//		}
//		err := k.kac.CreateTopic(topic, topicDetail, false)
//		if err != nil {
//			result.Err = err.Error()
//			return result
//		}
//	}
//
//	return result
//}

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

//
//// DescribeTopic 获取主题详情
//func (k *Service) DescribeTopic(topics []string) *types.ResultsResp {
//	result := &types.ResultsResp{}
//	if k.kac == nil {
//		result.Err = "请先选择连接"
//		return result
//	}
//	ctx := context.Background()
//
//	for _, topic := range topics {
//		topicMetadata, err := k.kac.DescribeTopicConfigs(ctx, topics...)
//		if err != nil {
//			result.Err = err.Error()
//			return result
//		}
//
//		for _, metadata := range topicMetadata {
//			var partitions []map[string]interface{}
//			for _, partition := range metadata.Partitions {
//				partitions = append(partitions, map[string]interface{}{
//					"Version":         partition.Version,
//					"Err":             partition.Err.Error(),
//					"ID":              partition.ID,
//					"Leader":          partition.Leader,
//					"LeaderEpoch":     partition.LeaderEpoch,
//					"Replicas":        partition.Replicas,
//					"Isr":             partition.Isr,
//					"OfflineReplicas": partition.OfflineReplicas,
//				})
//			}
//
//			result.Results = append(result.Results, map[string]interface{}{
//				"topic":      topic,
//				"Name":       metadata.Name,
//				"Err":        metadata.Err,
//				"Partitions": partitions,
//				"Version":    metadata.Version,
//				"IsInternal": metadata.IsInternal,
//				"Uuid":       metadata.Uuid,
//			})
//		}
//	}
//
//	return result
//}
//
//// CreatePartitions 添加分区
//func (k *Service) CreatePartitions(topics []string, count int32) *types.ResultResp {
//	result := &types.ResultResp{}
//	if k.kac == nil {
//		result.Err = "请先选择连接"
//		return result
//	}
//	for _, topic := range topics {
//		err := k.kac.CreatePartitions(topic, count, nil, false)
//		if err != nil {
//			result.Err = err.Error()
//			return result
//		}
//	}
//
//	return result
//}

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

//
//// AlterTopicConfig 修改主题配置
//func (k *Service) AlterTopicConfig(topic string, configs map[string]*string) *types.ResultsResp {
//	result := &types.ResultsResp{}
//	if k.kac == nil {
//		result.Err = "请先选择连接"
//		return result
//	}
//
//	err := k.kac.AlterConfig(sarama.TopicResource, topic, configs, false)
//	if err != nil {
//		result.Err = err.Error()
//		return result
//	}
//	return result
//}
//
//// DescribeGroup 消费组详情
//func (k *Service) DescribeGroup(groupID string) (*types.GroupInfo, error) {
//	// 获取消费组详情
//	groups, err := k.kac.DescribeConsumerGroups([]string{groupID})
//	if err != nil {
//		return nil, fmt.Errorf("describe consumer group failed: %v", err)
//	}
//	if len(groups) == 0 {
//		return nil, fmt.Errorf("group not found: %s", groupID)
//	}
//
//	// 获取消费组offset信息
//	offsetFetch, err := k.kac.ListConsumerGroupOffsets(groupID, nil)
//	if err != nil {
//		return nil, fmt.Errorf("list consumer group offsets failed: %v", err)
//	}
//
//	info := &types.GroupInfo{
//		Group:  groupID,
//		Topics: make(map[string][]types.PartitionOffset),
//	}
//
//	// 遍历每个topic的分区offset
//	for topic, partitions := range offsetFetch.Blocks {
//		var partitionOffsets []types.PartitionOffset
//
//		// 获取topic的最新offset
//		latestOffsets, err := k.getTopicLatestOffsets(topic)
//		if err != nil {
//			return nil, err
//		}
//
//		for partition, offsetBlock := range partitions {
//			if offsetBlock.Offset == -1 {
//				continue // 跳过未消费的分区
//			}
//
//			latestOffset := latestOffsets[partition]
//			lag := latestOffset - offsetBlock.Offset
//
//			po := types.PartitionOffset{
//				Partition: partition,
//				Offset:    offsetBlock.Offset,
//				Lag:       lag,
//			}
//			partitionOffsets = append(partitionOffsets, po)
//			info.TotalLag += lag
//		}
//		info.Topics[topic] = partitionOffsets
//	}
//
//	return info, nil
//}
//
//// 获取Topic各分区最新Offset
//func (k *Service) getTopicLatestOffsets(topic string) (map[int32]int64, error) {
//	// 获取topic的所有分区
//	metadata, err := k.kac.DescribeTopics([]string{topic})
//	if err != nil {
//		return nil, fmt.Errorf("describe topic failed: %v", err)
//	}
//
//	if len(metadata) == 0 {
//		return nil, fmt.Errorf("topic not found: %s", topic)
//	}
//
//	client, err := sarama.NewClient(k.bootstrapServers, k.config)
//	if err != nil {
//		return nil, fmt.Errorf("create client failed: %v", err)
//	}
//	defer client.Close()
//
//	offsets := make(map[int32]int64)
//	for _, partition := range metadata[0].Partitions {
//		offset, err := client.GetOffset(topic, partition.ID, sarama.OffsetNewest)
//		if err != nil {
//			return nil, fmt.Errorf("get partition offset failed: %v", err)
//		}
//		offsets[partition.ID] = offset
//	}
//
//	return offsets, nil
//}

//
//func (k *Service) GetTopicOffsets(topics []string, groupID string) *types.ResultResp {
//	result := &types.ResultResp{}
//
//	if k.kac == nil {
//		result.Err = "请先选择连接"
//		return result
//	}
//
//	topicOffset := make(map[string]map[int32][]int64)
//	topicLag := make(map[string][]int64)
//	group, err := sarama.NewConsumerGroup(
//		k.bootstrapServers,
//		groupID,
//		k.config,
//	)
//	sarama.NewConsumer()
//	if err != nil {
//		result.Err = err.Error()
//		return result
//	}
//
//	for _, topic := range topics {
//		partitions, err := group.Partitions(topic)
//		if err != nil {
//			continue
//		}
//
//		topicOffset[topic] = make(map[int32][]int64)
//		var totalEndOffset, totalCommitted int64
//
//		for _, partition := range partitions {
//			pc, err := group.ConsumePartition(topic, partition, sarama.OffsetOldest)
//			if err != nil {
//				continue
//			}
//
//			committed := pc.HighWaterMarkOffset()
//			endOffset, err := k.getLatestOffset(topic, partition)
//			if err != nil {
//				continue
//			}
//
//			lag := endOffset - committed
//			topicOffset[topic][partition] = []int64{committed, endOffset, lag}
//
//			totalEndOffset += endOffset
//			totalCommitted += committed
//		}
//
//		topicLag[topic] = []int64{totalEndOffset, totalCommitted}
//	}
//	result.Result = map[string]interface{}{"offsets": topicOffset, "lag": topicLag}
//
//	return result
//}

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
