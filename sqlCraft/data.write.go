package sqlCraft

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
	Payload []any `json:"payload"` // 消息载荷
}
