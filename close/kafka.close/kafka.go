package kafka_close

import (
	"github.com/IBM/sarama"
	"gitlab.itestor.com/helei/utils.git/ulogs"
)

func CloseKafkaPartition(partition sarama.PartitionConsumer) {
	if partition == nil {
		return
	}
	ulogs.Checkerr(partition.Close(), "CloseKafkaPartition")
}
