package kafkaHander

import (
	"fmt"
	"github.com/IBM/sarama"
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
				return fmt.Errorf("kafka状态失败:topic:%s,partition:%d", claim.Topic(), claim.Partition())
			}
			h.msg <- message
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
