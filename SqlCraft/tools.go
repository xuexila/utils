package SqlCraft

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
