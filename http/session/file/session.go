package file

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/dataType"
	"github.com/helays/utils/http/session"
	"github.com/helays/utils/tools"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
// Date: 2024/12/8 1:50
//

//var (
//	db *badger.DB
//)

// Instance session 实例
type Instance struct {
	option *session.Options
	Path   string `json:"path" yaml:"path" ini:"path"` // db路径
	ctx    context.Context
	cancel context.CancelFunc
}

// New 初始化 session 内存 实例
func New(opt ...Instance) (*Instance, error) {
	ins := &Instance{
		Path: "runtime/session",
	}
	if len(opt) > 0 {
		ins.Path = opt[0].Path
	}
	ins.Path = tools.Fileabs(ins.Path)
	err := tools.Mkdir(ins.Path)
	if err != nil {
		return nil, fmt.Errorf("创建session文件存放目录失败")
	}

	return ins, nil
}

// Register 注册结构定义
// 在使用文件作为session引擎的时候，需要将存储session值的结构注册进来。
func (this *Instance) Register(value ...any) {
	if len(value) < 1 {
		return
	}
	for _, v := range value {
		gob.Register(v)
	}
}

// Apply 应用配置
func (this *Instance) Apply(options *session.Options) {
	this.option = options
	this.ctx, this.cancel = context.WithCancel(context.Background())
	tools.RunAsyncTickerProbabilityFunc(this.ctx, !this.option.DisableGc, this.option.CheckInterval, this.option.GcProbability, this.gc)
}

// gc 垃圾回收
func (this *Instance) gc() {
	files, err := os.ReadDir(this.Path)
	if err != nil {
		return
	}
	// 循环所有文件
	// 如果是文件夹，就直接删除
	// 如果文件打开失败，跳过处理
	// 如果解析失败，就删除
	// 判断是否过期，过期也直接删除
	for _, file := range files {
		// 读取所有session文件
		sessionPath := filepath.Join(this.Path, file.Name())
		if file.IsDir() {
			_ = os.RemoveAll(sessionPath)
			continue
		}
		sessionVal := &session.Session{}

		f, err := os.Open(sessionPath)
		if err != nil {
			vclose.Close(f)
			continue
		}
		if err = gob.NewDecoder(f).Decode(sessionVal); err != nil {
			vclose.Close(f)
			_ = os.RemoveAll(sessionPath)
		}
		vclose.Close(f)
		if time.Time(sessionVal.ExpireTime).Before(time.Now()) {
			_ = os.RemoveAll(sessionPath)
		}
	}
}

// Close 关闭 db
func (this *Instance) Close() error {
	this.cancel()
	return nil
}

// 从session 文件中读取session 数据
func (this *Instance) get(w http.ResponseWriter, r *http.Request, name string) (*session.Session, string, error) {
	sessionId, err := session.GetSessionId(w, r, this.option) // 这一步一般不会失败
	if err != nil {
		return nil, "", err // 从cookie中获取sessionId失败
	}
	_k := session.GetSessionName(sessionId, name)
	// 从文件中读取数据
	sessionPath := filepath.Join(this.Path, _k)
	f, err := os.Open(sessionPath)
	defer vclose.Close(f)
	if err != nil {
		return nil, "", err
	}
	sessionVal := &session.Session{}
	if err = gob.NewDecoder(f).Decode(sessionVal); err != nil {
		vclose.Close(f)
		// session 数据解析失败，删除session文件
		_ = os.Remove(sessionPath)
		return nil, "", err
	}
	vclose.Close(f)
	if time.Time(sessionVal.ExpireTime).Before(time.Now()) {
		_ = os.Remove(sessionPath)
		// session已过期
		return nil, "", session.ErrNotFound
	}
	return sessionVal, _k, nil
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
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	return nil
}

// GetUp 获取session并更新过期时间
func (this *Instance) GetUp(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	// 更新session过期时间
	sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
	if err = this.set(w, r, *sessionVal); err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
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
	_ = os.Remove(filepath.Join(this.Path, _k))
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	return nil
}

// 设置session，将其通过gob 写入文件
func (this *Instance) set(w http.ResponseWriter, r *http.Request, dst session.Session) error {
	_k := session.GetSessionName(dst.Id, dst.Name)
	sessionPath := filepath.Join(this.Path, _k)
	f, err := os.OpenFile(sessionPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer vclose.Close(f)
	if err != nil {
		return err
	}
	return gob.NewEncoder(f).Encode(dst)
}

// Set 设置session
// w
// r
// name  session 名称
// value session 值
// duration session 过期时间，默认为24小时
func (this *Instance) Set(w http.ResponseWriter, r *http.Request, name string, value any, duration ...time.Duration) error {
	sessionId, _ := session.GetSessionId(w, r, this.option)
	sessionVal := session.Session{
		Id:         sessionId,
		Name:       name,
		Values:     session.SessionValue{Val: value},
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
	return this.set(w, r, sessionVal)
}

// Del 删除session
func (this *Instance) Del(w http.ResponseWriter, r *http.Request, name string) error {
	sessionId, _ := session.GetSessionId(w, r, this.option)
	_k := session.GetSessionName(sessionId, name)
	_ = os.Remove(filepath.Join(this.Path, _k))
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
	files, err := os.ReadDir(this.Path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), sessionId) {
			filePath := filepath.Join(this.Path, file.Name())
			_ = os.RemoveAll(filePath)
		}
	}
	return nil
}
