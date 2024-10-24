package kafka_close

import (
	"github.com/IBM/sarama"
	"github.com/xuexila/utils/ulogs"
)

func CloseKafkaPartition(partition sarama.PartitionConsumer) {
	if partition == nil {
		return
	}
	ulogs.Checkerr(partition.Close(), "CloseKafkaPartition")
}
