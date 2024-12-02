package sqlCraft

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

func (this SqlFilter) Builder(inputTx *gorm.DB) (*gorm.DB, error) {
	tx := this.TableName(inputTx)
	this.SetQuota(inputTx)
	this._select(tx)
	if err := this._join(tx); err != nil {
		return nil, err
	}
	this.SetWhere(tx)
	if this.Count {
		return tx, nil
	}
	this._group(tx)
	this._order(tx)
	this._offset(tx)
	this._limit(tx)
	return tx, nil
}

var (
	fieldQueryRegexp = regexp.MustCompile(`^[\w)( +-/*,:.=<>']+$`)
	seriesNoRegexp   = regexp.MustCompile(`^[a-z0-9_]{24}$`)
)

func (this SqlFilter) _select(_tx *gorm.DB) {
	var selectField []string
	for _, item := range this.Select {
		if fieldQueryRegexp.MatchString(item) { // 通过正则校验按
			//holder = append(holder, "?")
			selectField = append(selectField, fieldAliasAddQuota(item, this.quota))
		}
	}
	if len(selectField) > 0 {
		_tx.Select(selectField)
	}
}

func (this SqlFilter) _join(tx *gorm.DB) error {
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
		join[2] = fieldAliasAddQuota(join[2], this.quota)
		var _join string
		_join = fmt.Sprintf(`%s %s%s%s %s%s%s on %s`,
			strings.ToLower(join[0]),
			this.quota,
			joinTable,
			this.quota,
			this.quota,
			join[1],
			this.quota,
			join[2],
		)
		tx.Joins(_join)
	}
	return nil
}

func (this SqlFilter) SetWhere(tx *gorm.DB) {
	sqlStr, valParam := this.Where.sql(this.quota)
	if sqlStr == "" {
		return
	}
	tx.Where(sqlStr, valParam...)
}

func (this SqlFilter) _group(tx *gorm.DB) {
	var selectField []string
	for _, item := range this.Group {
		if fieldQueryRegexp.MatchString(item) { // 通过正则校验按
			selectField = append(selectField, fieldAliasAddQuota(item, this.quota))
		}
	}
	if len(selectField) > 0 {
		tx.Group(strings.Join(selectField, ","))
	}
}

func (this SqlFilter) _order(tx *gorm.DB) {
	for _, item := range this.Order {
		if fieldQueryRegexp.MatchString(item) {
			tx.Order(fieldAliasAddQuota(item, this.quota))
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
