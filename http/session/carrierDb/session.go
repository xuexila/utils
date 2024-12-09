package carrierDb

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/helays/utils/dataType"
	"github.com/helays/utils/db/userDb"
	"github.com/helays/utils/http/session"
	"github.com/helays/utils/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"reflect"
	"time"
)

type Instance struct {
	option *session.Options
	db     *gorm.DB
	ctx    context.Context
	cancel context.CancelFunc
}

// New 创建一个session实例
func New(db *gorm.DB, tableName ...string) *Instance {
	ins := &Instance{
		db: db,
	}
	userDb.AutoCreateTableWithStruct(db, session.Session{}, "自动创建session数据存储表失败")
	return ins
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

// Close 关闭 db
func (this *Instance) Close() error {
	this.cancel()
	return nil
}

// 自动gc
func (this *Instance) gc() {
	// 清理所有 过期时间不大于当前时间的数据
	this.db.Where("expire_time <= ?", time.Now()).Delete(&session.Session{})
}

// 从数据库中查询session 数据
func (this *Instance) get(w http.ResponseWriter, r *http.Request, name string) (sessionVal *session.Session, sessionId string, err error) {
	sessionId, err = session.GetSessionId(w, r, this.option) // 这一步一般不会失败
	if err != nil {
		return nil, "", err // 从cookie中获取sessionId失败
	}
	// 这里直接使用sessionId 和 name 去数据库查询
	tx := this.db.Where("id=? and name =?", sessionId, name).Take(sessionVal)
	if err = tx.Error; err != nil {
		return nil, "", err
	}
	return sessionVal, sessionId, nil
}

// 存储session数据
func (this *Instance) set(w http.ResponseWriter, r *http.Request, dst *session.Session) error {
	return this.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "name"}},
		UpdateAll: true,
	}).Create(dst).Error
}

// 删除session数据
func (this *Instance) del(sessionId, name string) {
	this.db.Where("id=? and name=?", sessionId, name).Delete(&session.Session{})
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
	if err = this.set(w, r, sessionVal); err != nil {
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
	sessionVal, sessionId, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	this.del(sessionId, name)
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
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
	now := time.Now()
	sessionVal := session.Session{
		Id:         sessionId,
		Name:       name,
		Values:     session.SessionValue{Val: value},
		CreateTime: dataType.CustomTime(now),
		Duration:   session.ExpireTime,
	}
	if len(duration) > 0 {
		sessionVal.Duration = duration[0]
	}
	sessionVal.ExpireTime = dataType.CustomTime(now.Add(sessionVal.Duration)) // 设置过期时间
	return this.set(w, r, &sessionVal)
}

// Del 删除session
func (this *Instance) Del(w http.ResponseWriter, r *http.Request, name string) error {
	sessionId, _ := session.GetSessionId(w, r, this.option)
	_k := session.GetSessionName(sessionId, name)
	this.del(sessionId, _k)
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
	return this.db.Where("id=?", sessionId).Delete(&session.Session{}).Error
}
