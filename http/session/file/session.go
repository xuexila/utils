package file

import (
	"context"
	"github.com/dgraph-io/badger/v4"
	"github.com/helays/utils/http/session/sessionConfig"
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

var (
	db *badger.DB
)

// Instance session 实例
type Instance struct {
	option *sessionConfig.Options
	Path   string `json:"path" yaml:"path" ini:"path"` // db路径
	ctx    context.Context
	cancel context.CancelFunc
}

// New 初始化 session 内存 实例
func New(opt ...Instance) (*Instance, error) {
	ins := &Instance{
		Path: "db/session",
	}
	if len(opt) > 0 {
		ins.Path = opt[0].Path
	}
	var err error
	c := badger.DefaultOptions(ins.Path)
	db, err = badger.Open(c)
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (this *Instance) Apply(options *sessionConfig.Options) {
	this.option = options
	this.ctx, this.cancel = context.WithCancel(context.Background())

}

// Close 关闭 db
func (this *Instance) Close() error {
	this.cancel()
	if db != nil {
		return db.Close()
	}
	return nil
}
