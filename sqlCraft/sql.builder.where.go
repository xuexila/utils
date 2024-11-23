package sqlCraft

import (
	"fmt"
	"reflect"
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
// Date: 2024/11/23 15:05
//

// 生成一组查询
func (this whereStruct) sql(quota string) (string, []any) {
	var valParams []any
	if len(this.Conditions) > 0 {
		var (
			_sqlArr []string
			_type   = strings.ToLower(this.Type)
		)
		for _, item := range this.Conditions {
			_sqlStr, _v := item.sql(quota)
			if _sqlStr == "" {
				continue
			}
			_sqlArr = append(_sqlArr, _sqlStr)
			valParams = append(valParams, _v...)
		}
		if len(_sqlArr) == 0 {
			return "", nil
		}
		if len(_sqlArr) == 1 {
			return _sqlArr[0], valParams
		}
		return fmt.Sprintf("(%s)", strings.Join(_sqlArr, fmt.Sprintf(" %s ", _type))), valParams
	}
	if !fieldQueryRegexp.MatchString(this.Field) {
		return "", nil
	}
	this.Field = fieldAliasAddQuota(this.Field, quota)
	switch this.Operator {
	case "=", ">", "<", ">=", "<=", "<>", "like", "not like":
		return fmt.Sprintf(`%s %s ?`, this.Field, this.Operator), []any{this.Value}
	case "in", "not in":
		rtyp := reflect.TypeOf(this.Value)
		_v := this.Value
		// 参数如果不是 list ,就需要手动设置成list
		if rtyp.Kind() != reflect.Array && rtyp.Kind() != reflect.Slice {
			_v = []any{this.Value}
		}
		return fmt.Sprintf(`%s %s ?`, this.Field, this.Operator), []any{_v}
	case "null":
		return fmt.Sprintf(`%s IS NULL`, this.Field), nil
	case "notnull":
		return fmt.Sprintf(`%s IS NOT NULL`, this.Field), nil
	}
	return "", nil
}
