package memory

import (
	"context"
	"fmt"
	"github.com/helays/utils/dataType"
	"github.com/helays/utils/http/session"
	"github.com/helays/utils/tools"
	"net/http"
	"reflect"
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
	option *session.Options
	ctx    context.Context
	cancel context.CancelFunc
}

// New 初始化 session 内存 实例
func New() *Instance {
	return &Instance{}
}

func (this *Instance) Register(value ...any) {

}

func (this *Instance) Apply(options *session.Options) {
	this.option = options
	this.ctx, this.cancel = context.WithCancel(context.Background())
	// 还需要自动删除
	tools.RunAsyncTickerFunc(this.ctx, true, this.option.CheckInterval, func() {
		sessionStorage.Range(func(key, value any) bool {
			ss := value.(session.Session)
			if time.Time(ss.ExpireTime).Before(time.Now()) {
				sessionStorage.Delete(key)
			}
			return true
		})
	})
}

func (this *Instance) Close() error {
	this.cancel() // 关闭定时器
	return nil
}

var (
	sessionStorage sync.Map // 存储session
)

func (this *Instance) get(w http.ResponseWriter, r *http.Request, name string) (*session.Session, string, error) {
	sessionId, err := session.GetSessionId(w, r, this.option) // 这一步一般不会失败
	if err != nil {
		return nil, "", err // 从cookie中获取sessionId失败
	}
	_k := session.GetSessionName(sessionId, name)
	val, ok := sessionStorage.Load(_k)
	if !ok {
		return nil, "", session.ErrNotFound
	}
	sessionVal := val.(session.Session)
	if time.Time(sessionVal.ExpireTime).Before(time.Now()) {
		sessionStorage.Delete(_k)
		// session已过期
		return nil, "", session.ErrNotFound
	}
	return &sessionVal, _k, nil
}

// Get 获取session
func (this *Instance) Get(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values))
	return nil
}

// GetUp 获取session并更新过期时间
func (this *Instance) GetUp(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _k, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	// 更新session过期时间
	sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
	sessionStorage.Store(_k, *sessionVal)
	v.Elem().Set(reflect.ValueOf(sessionVal.Values))
	return nil
}

// Flashes 获取并删除session
func (this *Instance) Flashes(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _k, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	sessionStorage.Delete(_k)
	v.Elem().Set(reflect.ValueOf(sessionVal.Values))
	return nil
}

// Set 设置session
// w
// r
// name  session 名称
// value session 值
// duration session 过期时间，默认为24小时
func (this *Instance) Set(w http.ResponseWriter, r *http.Request, name string, value any, duration ...time.Duration) error {
	sessionId, _ := session.GetSessionId(w, r, this.option)
	_k := session.GetSessionName(sessionId, name)
	sessionVal := session.Session{
		Id:         sessionId,
		Name:       name,
		Values:     value,
		CreateTime: dataType.CustomTime{},
		ExpireTime: dataType.CustomTime{},
		Duration:   session.ExpireTime,
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
	sessionId, _ := session.GetSessionId(w, r, this.option)
	_k := session.GetSessionName(sessionId, name)
	sessionStorage.Delete(_k)
	return nil
}

// Destroy 销毁session
func (this *Instance) Destroy(w http.ResponseWriter, r *http.Request) error {
	sessionId, err := session.GetSessionId(w, r, this.option)
	if err != nil {
		return err // 从cookie中获取sessionId失败
	}
	// 需要删除 cookie 或者 header
	session.DeleteSessionId(w, this.option)
	// 删除所有以 sessionId 为前缀的 key
	sessionStorage.Range(func(key, value any) bool {
		sessVal := value.(session.Session)
		if sessVal.Id == sessionId {
			sessionStorage.Delete(key)
		}
		return true
	})
	return nil
}
