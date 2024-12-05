package sqlCraft

import (
	"fmt"
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
// Date: 2024/11/23 14:50
//

var (
	fieldAliasRegex  = regexp.MustCompile(`(\w+?)\.`)
	reservedFieldMap = map[string]bool{
		"GROUP": true,
	}
)

// 如果有 。号，。号前面的字符串需要加双引号
func fieldAliasAddQuota(field string, arg string) string {

	field = fieldAliasRegex.ReplaceAllString(field, fmt.Sprintf(`%s$1%s.`, arg, arg))
	if reservedFieldMap[strings.ToUpper(field)] {
		field = fmt.Sprintf(`%s%s%s`, arg, field, arg)
	}
	return field
}

//// 新建一个db
//func newDb(inputTx *gorm.DB, schema ...string) (*gorm.DB, error) {
//	newSession := inputTx.Session(&gorm.Session{NewDB: true})
//	err := setSchema(newSession, schema...)
//	if err != nil {
//		return newSession, err
//	}
//	return newSession, nil
//}
//
//// 设置 模式
//func setSchema(tx *gorm.DB, schema ...string) error {
//	if len(schema) > 0 && schema[0] != "" {
//		switch tx.Dialector.Name() {
//		case "postgres":
//			return tx.Exec("set search_path To ?;", clause.Table{
//				Name: schema[0],
//				Raw:  false,
//			}).Error
//		}
//	}
//	return nil
//}
