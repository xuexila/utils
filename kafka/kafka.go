package kafka

import (
	"github.com/IBM/sarama"
	"gitlab.itestor.com/helei/utils.git"
)

func CloseKafkaPartition(partition sarama.PartitionConsumer) {
	if partition == nil {
		return
	}
	utils.Checkerr(partition.Close(), "CloseKafkaPartition")
}
