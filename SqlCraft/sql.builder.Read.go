package SqlCraft

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strings"
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
// Date: 2024/11/23 11:38
//

type SqlFilter struct {
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
}

type whereStruct struct {
	Type       string        `json:"type"`
	Field      string        `json:"field"`
	Operator   string        `json:"operator"`
	Value      any           `json:"value"`
	Conditions []whereStruct `json:"conditions"`
}

func (this SqlFilter) Builder(inputTx *gorm.DB) (*gorm.DB, error) {
	tx, err := this.TableName(inputTx)
	if err != nil {
		return nil, err
	}
	quota := ""
	switch tx.Dialector.Name() {
	case "mysql":
		quota = "`"
	case "postgres":
		quota = `"`
	}
	this._select(tx, quota)
	if err = this._join(tx, quota); err != nil {
		return nil, err
	}
	this.SetWhere(tx, quota)
	if this.Count {
		return tx, nil
	}
	this._group(tx, quota)
	this._order(tx, quota)
	this._offset(tx)
	this._limit(tx)
	return tx, nil
}

var (
	fieldQueryRegexp = regexp.MustCompile(`^[\w)( +-/*,:.=<>']+$`)
	seriesNoRegexp   = regexp.MustCompile(`^[a-z0-9_]{24}$`)
)

func (this SqlFilter) _select(_tx *gorm.DB, quota string) {
	var selectField []string
	for _, item := range this.Select {
		if fieldQueryRegexp.MatchString(item) { // 通过正则校验按
			//holder = append(holder, "?")
			selectField = append(selectField, fieldAliasAddQuota(item, quota))
		}
	}
	if len(selectField) > 0 {
		_tx.Select(selectField)
	}
}

func (this SqlFilter) _join(tx *gorm.DB, quota string) error {
	if len(this.Join) < 1 {
		return nil
	}
	for _, join := range this.Join {
		switch strings.ToLower(join[0]) {
		case "inner join", "left join", "right join", "join":
		default:
			return errors.New("不支持的连表方式：" + join[0])
		}
		if !fieldQueryRegexp.MatchString(join[2]) {
			return errors.New(fmt.Sprintf("连表查询校验不通过：%s", strings.Join(join[:], " ")))
		}
		// join on 里面的字符串加上引号

		joinTable := join[1]
		if this.FindTable != nil {
			tb, err := this.FindTable(join[1])
			if err != nil {
				return err
			}
			joinTable = tb
		}
		join[2] = fieldAliasAddQuota(join[2], quota)
		var _join string
		_join = fmt.Sprintf(`%s %s%s%s %s%s%s on %s`,
			strings.ToLower(join[0]),
			quota,
			joinTable,
			quota,
			quota,
			join[1],
			quota,
			join[2],
		)
		tx.Joins(_join)
	}
	return nil
}

func (this SqlFilter) SetWhere(tx *gorm.DB, quota string) {

	sqlStr, valParam := this.Where.sql(quota)
	if sqlStr == "" {
		return
	}

	tx.Where(sqlStr, valParam...)
}

func (this SqlFilter) _group(tx *gorm.DB, quota string) {
	var selectField []string
	for _, item := range this.Group {
		if fieldQueryRegexp.MatchString(item) { // 通过正则校验按
			selectField = append(selectField, fieldAliasAddQuota(item, quota))
		}
	}
	if len(selectField) > 0 {
		tx.Group(strings.Join(selectField, ","))
	}
}

func (this SqlFilter) _order(tx *gorm.DB, quota string) {
	for _, item := range this.Order {
		if fieldQueryRegexp.MatchString(item) {
			tx.Order(fieldAliasAddQuota(item, quota))

		}
	}
}

func (this SqlFilter) _offset(tx *gorm.DB) {
	if this.Offset > 0 {
		tx.Offset(this.Offset)
	}
}

func (this SqlFilter) _limit(tx *gorm.DB) {
	if this.Limit > 0 {
		tx.Limit(this.Limit)
	} else if this.Limit == 0 {
		tx.Limit(1000)
	}
}
