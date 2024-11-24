package sqlCraft

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/helays/utils/http/httpServer"
	"github.com/helays/utils/tools"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2024/11/24 18:07
//

// DataWrite 数据写
type DataWrite struct {
	Debug  bool   `json:"-"`      // 是否调试模式
	Remark string `json:"remark"` // 数据库唯一标识
	Mode   string `json:"mode"`   // sync 同步模式 async 异步模式
	// 数据库相关
	Table  string `json:"table"`  // 表名
	Schema string `json:"schema"` // 表模式
	// topic
	Topic  string            `json:"topic"`  // topic
	Key    string            `json:"key"`    // 消息key
	Header map[string]string `json:"header"` // 消息header

	//
	Payload any `json:"payload"` // 消息载荷
}

func (this *DataWrite) Save2Db(w http.ResponseWriter, r *http.Request, inputTx *gorm.DB) {
	if this.Payload == nil {
		httpServer.SetReturnCode(w, r, 500, "无有效载荷")
		return
	}
	var saveData []map[string]any
	// 判断 Payload是否是数组，如果不是数组，需要转换为数组
	if _d, ok := this.Payload.(map[string]any); ok {
		saveData = append(saveData, _d)
	} else if _, ok = this.Payload.([]any); ok {
		for _, v := range this.Payload.([]any) {
			if _d, ok = v.(map[string]any); ok {
				saveData = append(saveData, _d)
			}
		}
	} else {
		httpServer.SetReturnCode(w, r, 500, "载荷无法识别入库", this.Payload)
		return
	}
	if _, ok := this.Payload.([]any); !ok {
		this.Payload = []any{this.Payload}
	}
	uTx, err := newDb(inputTx, this.Schema)
	if err != nil {
		httpServer.SetReturnError(w, r, err, 500, "数据库复制失败")
		return
	}
	err = uTx.Table(this.Table).CreateInBatches(&saveData, 1000).Error
	if err != nil {
		httpServer.SetReturnError(w, r, err, 500, "数据写入失败")
		return
	}
	httpServer.SetReturnCode(w, r, 0, "数据写入成功", uTx.RowsAffected)
}

func (this *DataWrite) Save2Kafka(w http.ResponseWriter, r *http.Request, sd func(msg *sarama.ProducerMessage) error) {
	// 判断 Payload是否是数组，如果不是数组，需要转换为数组
	if _, ok := this.Payload.([]any); !ok {
		this.Payload = []any{this.Payload}
	}
	var errs []string
	for idx, mesage := range this.Payload.([]any) {
		// 构造消息
		byt, err := tools.Any2bytes(mesage)
		if err != nil {
			errs = append(errs, fmt.Sprintf("第%d条数据转换失败: %s；%v", idx, err, mesage))
			continue
		}
		msg := &sarama.ProducerMessage{
			Topic: this.Topic,
			Value: sarama.ByteEncoder(byt),
		}
		if this.Key != "" {
			msg.Key = sarama.StringEncoder(this.Key)
		}
		for k, v := range this.Header {
			msg.Headers = append(msg.Headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
		err = sd(msg)
		if err != nil {
			errs = append(errs, fmt.Sprintf("第%d条数据写入失败: %s；%v", idx, err, mesage))
			continue
		}
	}
	if len(errs) < 1 {
		httpServer.SetReturnCode(w, r, 0, "数据写入成功")
		return
	}
	httpServer.SetReturnCode(w, r, 500, "数据写入失败", strings.Join(errs, "\n"))
}
