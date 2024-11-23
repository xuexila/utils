package sqlCraft

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

type SqlFilter struct {
	Sql          string                          `json:"sql"`
	Args         []any                           `json:"args"`
	Type         string                          `json:"type"` // query exec
	Select       []string                        `json:"select"`
	Count        bool                            `json:"count"`
	Table        string                          `json:"table"`
	Alias        string                          `json:"alias"`
	Where        whereStruct                     `json:"where"`
	Order        []string                        `json:"order"`
	Join         [][3]string                     `json:"join"`
	Limit        int                             `json:"limit"`
	Offset       int                             `json:"offset"`
	Group        []string                        `json:"group"`
	Set          any                             `json:"set"`
	Export       bool                            `json:"export"`        // 是否导出
	ExportHeader map[string]string               `json:"export_header"` // 导出自定义表头
	FileType     string                          `json:"file_type"`     // 导出文件格式，可选excel,csv。默认excel
	Schema       string                          `json:"schema"`        // 数据库模式
	Debug        bool                            `json:"-"`             // 是否调试模式
	Remark       string                          `json:"remark"`        // 数据库唯一标识
	FindTable    func(id string) (string, error) // 根据ID获取
	quota        string                          // 引号类型
}

type whereStruct struct {
	Type       string        `json:"type"`
	Field      string        `json:"field"`
	Operator   string        `json:"operator"`
	Value      any           `json:"value"`
	Conditions []whereStruct `json:"conditions"`
}

// NewDb 克隆 一个事务
func (this SqlFilter) NewDb(inputTx *gorm.DB) (*gorm.DB, error) {
	newSession := inputTx.Session(&gorm.Session{NewDB: true})
	err := this.setSchema(newSession)
	if err != nil {
		return newSession, err
	}
	return newSession, nil
}

// 设置 模式
func (this SqlFilter) setSchema(tx *gorm.DB) error {
	if this.Schema != "" {
		switch tx.Dialector.Name() {
		case "postgres":
			return tx.Exec("set search_path To ?;", clause.Table{
				Name: this.Schema,
				Raw:  false,
			}).Error
		}
	}
	return nil
}

// TableName 设置表名
func (this SqlFilter) TableName(inputTx *gorm.DB) *gorm.DB {
	if this.Debug {
		return inputTx.Debug().Table("?", clause.Table{
			Name:  this.Table,
			Alias: this.Alias,
			Raw:   false,
		})
	}
	return inputTx.Table("?", clause.Table{
		Name:  this.Table,
		Alias: this.Alias,
		Raw:   false,
	})
}

func (this *SqlFilter) SetQuota(tx *gorm.DB) {
	switch tx.Dialector.Name() {
	case "mysql":
		this.quota = "`"
	case "postgres":
		this.quota = `"`
	}
}
