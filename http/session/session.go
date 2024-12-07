package session

import (
	"github.com/helays/utils/http/session/sessionConfig"
	"github.com/helays/utils/tools"
	"time"
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
// Date: 2024/12/7 22:02
//

// session 支持内存存储、程序重启session失效
// 支持文件存储，可持久化存储
// 支持 redis 存储
// 支持数据库存储

type Store struct {
	*sessionConfig.Options
	sessionConfig.Dialector
}

// Init 初始化session
func Init(dialector sessionConfig.Dialector, opt ...*sessionConfig.Options) *Store {
	c := &Store{}
	if len(opt) > 0 {
		c.Options = opt[0]
		c.Options.CheckInterval = tools.AutoTimeDuration(c.Options.CheckInterval, time.Microsecond, sessionConfig.Interval)
	} else {
		c.Options.CheckInterval = sessionConfig.Interval
		c.Options.Carrier = "cookie"
	}
	c.Options.CookieName = tools.Ternary(c.Options.CookieName == "", sessionConfig.CookieName, c.Options.CookieName)
	c.Dialector = dialector
	c.Dialector.Apply(c.Options)
	return c
}
