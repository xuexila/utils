package kafkaClose

import (
	"github.com/IBM/sarama"
	"github.com/helays/utils/ulogs"
)

// CloseKafkaPartition 关闭分区
func CloseKafkaPartition(partition sarama.PartitionConsumer) {
	if partition == nil {
		return
	}
	ulogs.Checkerr(partition.Close(), "CloseKafkaPartition")
}

// CloseKafkaConsumerGroup 关闭消费者组
func CloseKafkaConsumerGroup(group sarama.ConsumerGroup) {
	if group != nil {
		ulogs.Checkerr(group.Close(), "CloseKafkaConsumerGroup 执行失败")
	}
}
