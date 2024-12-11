package test

import (
	"github.com/helays/utils/http/session"
	"github.com/helays/utils/http/session/carrierMemory"
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
// Date: 2024/12/8 17:07
//

var (
	store = &session.Store{}
)

func run() {
	defer session.Close(store)
	store = session.Init(carrierMemory.New(), &session.Options{
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
