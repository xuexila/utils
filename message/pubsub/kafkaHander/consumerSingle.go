package kafkaHander

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/helays/utils/message/pubsub"
	"time"
)

func (this *Instance) single(param pubsub.Params) {
	partitionList, err := this.consumer.Partitions(param.Topic)
	if err != nil {
		this.error("订阅发布组件订阅失败", "kafka consumer", fmt.Errorf("%s：%s", param.Topic, err.Error()))
		return
	}
	for _, partition := range partitionList {
		go this.partition(param.Topic, partition, sarama.OffsetNewest) // 默认从最新的offset开始消费
	}
}

// 分区消费
func (this *Instance) partition(topic string, partition int32, offset int64) {
	this.log("订阅发布组件", "kafka载体", "普通消费者", topic, "开始消费", "分区", partition)
	pc, err := this.consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		this.error("订阅发布组件订阅失败", "kafka consumer", fmt.Errorf("%s：%s", topic, err.Error()))
		return
	}
	defer pc.AsyncClose()
	for {
		select {
		case msg := <-pc.Messages(): // 收取数据
			if msg == nil {
				this.error("订阅发布组件订阅失败", "kafka consumer", topic, "消息为空", "分区", partition)
				continue
			}
			this.message <- msg
		case <-this.opts.Ctx.Done(): // 监听退出信号
			this.log("订阅发布组件", "kafka载体", "普通消费者", topic, "退出消费", "分区", partition)
			return
		case err = <-pc.Errors(): // 监听错误
			this.error("订阅发布组件订阅失败", "kafka consumer", fmt.Errorf("%s：%s", topic, err.Error()), "分区", partition)
			if errors.Is(err, sarama.ErrOffsetOutOfRange) || err == nil {
				pc.AsyncClose()
				pc, err = this.consumer.ConsumePartition(topic, partition, offset)
				if err != nil {
					time.Sleep(10 * time.Second)
					this.error("订阅发布组件订阅失败", "kafka consumer", "重新消费分区失败", fmt.Errorf("%s：%s", topic, err.Error()), "分区", partition)
					continue
				}
			}
		}
	}
}
