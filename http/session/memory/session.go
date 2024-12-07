package memory

import (
	"context"
	"github.com/helays/utils/dataType"
	"github.com/helays/utils/http/session/sessionConfig"
	"github.com/helays/utils/tools"
	"net/http"
	"strings"
	"sync"
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
// Date: 2024/12/7 23:23
//
// vsclub:sessionId:key:value

// Instance session 实例
type Instance struct {
	option *sessionConfig.Options
}

// New 初始化 session 内存 实例
func New() *Instance {
	return &Instance{}
}

func (this *Instance) Apply(options *sessionConfig.Options) {
	this.option = options
	// 还需要自动删除
	tools.RunAsyncTickerFunc(context.Background(), true, this.option.CheckInterval, func() {
		sessionStorage.Range(func(key, value any) bool {
			session := value.(sessionConfig.Session)
			if time.Time(session.ExpireTime).Before(time.Now()) {
				sessionStorage.Delete(key)
			}
			return true
		})
	})
}

var (
	sessionStorage sync.Map // 存储session
)

func (this *Instance) get(w http.ResponseWriter, r *http.Request, name string) (*sessionConfig.Session, string, error) {
	sessionId, err := sessionConfig.GetSessionId(w, r, this.option)
	if err != nil {
		return nil, "", err // 从cookie中获取sessionId失败
	}
	_k := sessionConfig.GetSessionName(sessionId, name)
	val, ok := sessionStorage.Load(_k)
	if !ok {
		return nil, "", sessionConfig.ErrNotFound
	}
	sesseionVal := val.(sessionConfig.Session)
	if time.Time(sesseionVal.ExpireTime).Before(time.Now()) {
		sessionStorage.Delete(_k)
		// session已过期
		return nil, "", sessionConfig.ErrNotFound
	}
	return &sesseionVal, _k, nil
}

// Get 获取session
func (this *Instance) Get(w http.ResponseWriter, r *http.Request, name string) (string, error) {
	sessionVal, _, err := this.get(w, r, name)
	if err != nil {
		return "", err
	}
	return sessionVal.Values, nil
}

// GetUp 获取session并更新过期时间
func (this *Instance) GetUp(w http.ResponseWriter, r *http.Request, name string) (string, error) {
	sessionVal, _k, err := this.get(w, r, name)
	if err != nil {
		return "", err
	}
	// 更新session过期时间
	sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
	sessionStorage.Store(_k, *sessionVal)
	return sessionVal.Values, nil
}

// Flashes 获取并删除session
func (this *Instance) Flashes(w http.ResponseWriter, r *http.Request, name string) (string, error) {
	sessionVal, _k, err := this.get(w, r, name)
	if err != nil {
		return "", err
	}
	sessionStorage.Delete(_k)
	return sessionVal.Values, nil
}

func (this *Instance) Set(w http.ResponseWriter, r *http.Request, name string, value any, duration ...time.Duration) error {
	sessionId, _ := sessionConfig.GetSessionId(w, r, this.option)
	dstVal, err := tools.Any2bytes(value)
	if err != nil {
		return err
	}
	_k := sessionConfig.GetSessionName(sessionId, name)
	sessionVal := sessionConfig.Session{
		Id:         sessionId,
		Name:       name,
		Values:     string(dstVal),
		CreateTime: dataType.CustomTime{},
		ExpireTime: dataType.CustomTime{},
		Duration:   sessionConfig.ExpireTime,
	}
	if len(duration) > 0 {
		sessionVal.Duration = duration[0]
	}
	now := time.Now()
	sessionVal.Id = sessionId                                                 // 设置sessionId
	sessionVal.CreateTime = dataType.CustomTime(now)                          // 设置创建时间
	sessionVal.ExpireTime = dataType.CustomTime(now.Add(sessionVal.Duration)) // 设置过期时间
	sessionStorage.Store(_k, sessionVal)
	return nil
}

// Del 删除session
func (this *Instance) Del(w http.ResponseWriter, r *http.Request, name string) error {
	sessionId, _ := sessionConfig.GetSessionId(w, r, this.option)
	_k := sessionConfig.GetSessionName(sessionId, name)
	sessionStorage.Delete(_k)
	return nil
}

// Destroy 销毁session
func (this *Instance) Destroy(w http.ResponseWriter, r *http.Request) error {
	sessionId, err := sessionConfig.GetSessionId(w, r, this.option)
	if err != nil {
		return err // 从cookie中获取sessionId失败
	}
	// 需要删除 cookie 或者 header
	sessionConfig.DeleteSessionId(w, this.option)
	// 删除所有以 sessionId 为前缀的 key
	sessionStorage.Range(func(key, value any) bool {
		_k := key.(string)
		// 判断_k 是否以 sessionId 开头
		if strings.HasPrefix(_k, sessionId) {
			sessionStorage.Delete(key)
		}
		return true
	})
	return nil
}
