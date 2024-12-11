package dbFiledType

type FieldType struct {
	Name        string `json:"name"`        // 字段类型
	FullName    string `json:"full_name"`   // 完整类型
	Extra       string `json:"extra"`       // 补充说明
	DstDb       string `json:"dst_db"`      // 目标数据库类型
	DefaultLen  int    `json:"defaultLen"`  // 默认长度
	DecimalCols int    `json:"decimalCols"` // 小数点位数
}
