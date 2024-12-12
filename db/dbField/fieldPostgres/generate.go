package fieldPostgres

import (
	"fmt"
	"github.com/helays/utils/db/dbField"
	"strings"
)

// 设置字段类型
func generateFieldType(field dbField.Field, fieldInstance dbField.FieldInstance) (_t string) {
	if !field.LenReadOnly && fieldInstance.Length > 0 {
		field.DefaultLen = fieldInstance.Length
	}
	if !field.ColsReadOnly && fieldInstance.Precision > 0 {
		field.DecimalCols = fieldInstance.Precision
	}

	_t = field.Name
	if field.DefaultLen == 0 && field.DecimalCols == 0 {
		return
	}
	if field.DecimalCols > 0 {
		_t += fmt.Sprintf("(%d,%d)", field.DefaultLen, field.DecimalCols)
		return
	}
	_t += fmt.Sprintf("(%d)", field.DefaultLen)
	return
}

// GenerateCreate 生成 create 字段
func GenerateCreate(field dbField.Field, fieldInstance dbField.FieldInstance) (string, []any) {
	var (
		sb        strings.Builder
		paramsVal []any
	)
	sb.WriteString(fmt.Sprintf(`"%s" `, fieldInstance.Field))
	// 设置字段类型
	sb.WriteString(generateFieldType(field, fieldInstance))
	// 是否允许 null
	if fieldInstance.Required {
		sb.WriteString(" not null")
	}
	if field.Attribute == "date" && fieldInstance.AutoCreateTime {
		sb.WriteString(" default current_timestamp")
	} else {
		_default := fieldInstance.Default
		if _default != "" {
			if _default == `''` || _default == `""` {
				sb.WriteString(" default ''")
			} else {
				sb.WriteString(" default ?")
				paramsVal = append(paramsVal, _default)
			}
		}
	}

	return sb.String() + ",", paramsVal
}

// GenerateAlter 生成 alter 字段
func GenerateAlter(field dbField.Field, fieldInstance dbField.FieldInstance) (string, []any) {
	var (
		alertStr  []string
		paramsVal []any
	)
	prefix := fmt.Sprintf(`ALTER COLUMN "%s"`, fieldInstance.Field)
	alertStr = append(alertStr, fmt.Sprintf("%s TYPE %s", prefix, generateFieldType(field, fieldInstance))) // alart 字段
	if fieldInstance.Required {
		alertStr = append(alertStr, fmt.Sprintf("%s SET NOT NULL", prefix))
	}
	if field.Attribute == "date" && fieldInstance.AutoCreateTime {
		alertStr = append(alertStr, fmt.Sprintf("%s SET DEFAULT current_timestamp", prefix))
	} else {
		_default := fieldInstance.Default
		if _default != "" {
			if _default == `''` || _default == `""` {
				alertStr = append(alertStr, fmt.Sprintf("%s SET DEFAULT ''", prefix))
			} else {
				alertStr = append(alertStr, fmt.Sprintf("%s SET DEFAULT ?", prefix))
				paramsVal = append(paramsVal, _default)
			}
		}
	}

	return strings.Join(alertStr, ",") + ",", paramsVal
}
