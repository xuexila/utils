package dbField

type Field struct {
	Name         string `json:"name"`           // 字段类型
	FullName     string `json:"full_name"`      // 完整类型
	Extra        string `json:"extra"`          // 补充说明
	DstDb        string `json:"dst_db"`         // 目标数据库类型
	LenReadOnly  bool   `json:"len_read_only"`  // 长度是否可写
	DefaultLen   int    `json:"default_len"`    // 默认长度
	ColsReadOnly bool   `json:"cols_read_only"` // 小数点位数是否可写
	DecimalCols  int    `json:"decimal_cols"`   // 小数点位数
	Attribute    string `json:"attribute"`      // 数字类型、字符串类型、日期类型等
}

type FieldInstance struct {
	Field             string `json:"field"`               // 字段
	Length            int    `json:"length"`              // 字段长度
	Precision         int    `json:"precision"`           // 字段精度
	Required          bool   `json:"required"`            // 是否必填， 3.0 新增
	Default           string `json:"default"`             // 默认值
	Primary           bool   `json:"primary"`             // 是否主键
	AutoIncrement     bool   `json:"auto_increment"`      // 是否自动增长
	AutoIncrementStep int    `json:"auto_increment_step"` // 自动增长步长
	AutoCreateTime    bool   `json:"auto_create_time"`    // 创建的时候自动写当前时间
	AutoUpdateTime    bool   `json:"auto_update_time"`    // 更新的时候是否自动写当前时间
	Extra             string `json:"extra"`               // 注释
}
