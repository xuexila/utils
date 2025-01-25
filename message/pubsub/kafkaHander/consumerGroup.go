package kafkaHander

import (
	"fmt"
	"github.com/helays/utils/message/pubsub"
)

func (this *Instance) group(param pubsub.Params) {
	var gh = consumerGroupHander{msg: this.message}
	this.log("订阅发布组件", "kafka consumerGroup", param.Topic)
	if err := this.consumerGroup.Consume(this.opts.Ctx, []string{param.Topic}, &gh); err != nil {
		this.error("订阅发布组件订阅失败", "kafka consumerGroup", fmt.Errorf("%s：%s", param.Topic, err.Error()))
	}
	// todo 这里需要 验证 group 订阅失败后的重连状态

}
