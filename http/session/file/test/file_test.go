package test

import (
	"github.com/helays/utils/http/session"
	"github.com/helays/utils/http/session/file"
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
// Date: 2024/12/8 17:11
//

var (
	store = &session.Store{}
)

type User struct{}

func run() {
	defer session.Close(store)
	engine, _ := file.New(file.Instance{Path: "runtime/session"})
	// 在session中需要存储User 结构体数据，需要将结构体注册进去
	// 需要在session 初始化之前进行注册
	engine.Register(User{})

	store = session.Init(engine, &session.Options{
		CookieName:    "vsclub.ltd",
		CheckInterval: time.Hour,
		Carrier:       "cookie",
		Path:          "",
		Domain:        "",
		MaxAge:        0,
		Secure:        false,
		HttpOnly:      false,
		SameSite:      0,
	})

}
