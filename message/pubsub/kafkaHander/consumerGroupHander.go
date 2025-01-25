package kafkaHander

import (
	"github.com/IBM/sarama"
	"github.com/helays/utils/logger/ulogs"
)

type consumerGroupHander struct {
	msg chan *sarama.ConsumerMessage
}

func (this *consumerGroupHander) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHander) Cleanup(session sarama.ConsumerGroupSession) error {
	// Optional: implement this if you need to clean up any state for the session.
	return nil
}

func (h *consumerGroupHander) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				ulogs.Log("kafka", "消费者关闭", claim.Topic(), claim.Partition())
				return nil
			}
			h.msg <- message

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
