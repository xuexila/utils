package SqlCraft

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
// Date: 2024/11/23 16:03
//

// 设置 模式
func (this SqlFilter) setSchema(tx *gorm.DB) error {
	var err error
	if this.Schema != "" {
		err = tx.Exec("set search_path To ?;", clause.Table{
			Name: this.Schema,
			Raw:  false,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// TableName 设置表名
func (this SqlFilter) TableName(inputTx *gorm.DB) (*gorm.DB, error) {
	newSession := inputTx.Session(&gorm.Session{NewDB: true})
	err := this.setSchema(newSession)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		return newSession.Debug().Table("?", clause.Table{
			Name:  this.Table,
			Alias: this.Alias,
			Raw:   false,
		}), nil
	}
	return newSession.Table("?", clause.Table{
		Name:  this.Table,
		Alias: this.Alias,
		Raw:   false,
	}), nil
}
