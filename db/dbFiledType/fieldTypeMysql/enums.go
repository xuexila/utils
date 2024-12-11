package fieldTypeMysql

import "github.com/helays/utils/db/dbFiledType"

var FieldTypeEnums = map[string]dbFiledType.FieldType{
	"tinyint": {
		Name:        "tinyint",
		FullName:    "tinyint",
		Extra:       "1字节有符号整数，范围-128到127；无符号时0到255",
		DstDb:       "mysql",
		DefaultLen:  4,
		DecimalCols: 0,
	},
	"smallint": {
		Name:        "smallint",
		FullName:    "smallint",
		Extra:       "2字节有符号整数，范围-32768到32767；无符号时0到65535",
		DstDb:       "mysql",
		DefaultLen:  6,
		DecimalCols: 0,
	},
	"mediumint": {
		Name:        "mediumint",
		FullName:    "mediumint",
		Extra:       "3字节有符号整数，范围-8388608到8388607；无符号时0到16777215",
		DstDb:       "mysql",
		DefaultLen:  9,
		DecimalCols: 0,
	},
	"int": {
		Name:        "int",
		FullName:    "int",
		Extra:       "4字节有符号整数，范围-2147483648到2147483647；无符号时0到4294967295",
		DstDb:       "mysql",
		DefaultLen:  11,
		DecimalCols: 0,
	},
	"integer": {
		Name:        "integer",
		FullName:    "integer",
		Extra:       "同int，4字节有符号整数",
		DstDb:       "mysql",
		DefaultLen:  11,
		DecimalCols: 0,
	},
	"bigint": {
		Name:        "bigint",
		FullName:    "bigint",
		Extra:       "8字节有符号整数，范围-9223372036854775808到9223372036854775807；无符号时0到18446744073709551615",
		DstDb:       "mysql",
		DefaultLen:  20,
		DecimalCols: 0,
	},
	"float": {
		Name:        "float",
		FullName:    "float",
		Extra:       "单精度浮点数，通常为4字节",
		DstDb:       "mysql",
		DefaultLen:  12,
		DecimalCols: 2,
	},
	"double": {
		Name:        "double",
		FullName:    "double",
		Extra:       "双精度浮点数，通常为8字节",
		DstDb:       "mysql",
		DefaultLen:  22,
		DecimalCols: 6,
	},
	"decimal": {
		Name:        "decimal",
		FullName:    "decimal",
		Extra:       "用户指定精度的准确小数",
		DstDb:       "mysql",
		DefaultLen:  10, // 默认总位数，具体取决于应用需求
		DecimalCols: 0,  // 需要明确指定小数点后位数，默认为0
	},
	"numeric": {
		Name:        "numeric",
		FullName:    "numeric",
		Extra:       "用户指定精度的准确小数 (与decimal相同)",
		DstDb:       "mysql",
		DefaultLen:  10, // 默认总位数，具体取决于应用需求
		DecimalCols: 0,  // 需要明确指定小数点后位数，默认为0
	},
	"date": {
		Name:        "date",
		FullName:    "date",
		Extra:       "存储日期 (年-月-日)",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"time": {
		Name:        "time",
		FullName:    "time",
		Extra:       "一天中的时间",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"datetime": {
		Name:        "datetime",
		FullName:    "datetime",
		Extra:       "存储日期和时间，没有时区信息",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"timestamp": {
		Name:        "timestamp",
		FullName:    "timestamp",
		Extra:       "存储UTC日期和时间，自动更新为当前时间戳",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"year": {
		Name:        "year",
		FullName:    "year",
		Extra:       "存储年份，范围1901到2155",
		DstDb:       "mysql",
		DefaultLen:  4,
		DecimalCols: 0,
	},
	"char": {
		Name:        "char",
		FullName:    "char",
		Extra:       "固定长度字符串，长度为n",
		DstDb:       "mysql",
		DefaultLen:  1,
		DecimalCols: 0,
	},
	"varchar": {
		Name:        "varchar",
		FullName:    "varchar",
		Extra:       "可变长度字符串，最大长度n",
		DstDb:       "mysql",
		DefaultLen:  255,
		DecimalCols: 0,
	},
	"tinytext": {
		Name:        "tinytext",
		FullName:    "tinytext",
		Extra:       "可变长度字符串，最大长度255字节",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"text": {
		Name:        "text",
		FullName:    "text",
		Extra:       "可变长度字符串，最大长度65535字节",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"mediumtext": {
		Name:        "mediumtext",
		FullName:    "mediumtext",
		Extra:       "可变长度字符串，最大长度16777215字节",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"longtext": {
		Name:        "longtext",
		FullName:    "longtext",
		Extra:       "可变长度字符串，最大长度4294967295字节",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"binary": {
		Name:        "binary",
		FullName:    "binary",
		Extra:       "固定长度二进制字符串，长度为n",
		DstDb:       "mysql",
		DefaultLen:  1,
		DecimalCols: 0,
	},
	"varbinary": {
		Name:        "varbinary",
		FullName:    "varbinary",
		Extra:       "可变长度二进制字符串，最大长度n",
		DstDb:       "mysql",
		DefaultLen:  255,
		DecimalCols: 0,
	},
	"tinyblob": {
		Name:        "tinyblob",
		FullName:    "tinyblob",
		Extra:       "二进制大对象，最大长度255字节",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"blob": {
		Name:        "blob",
		FullName:    "blob",
		Extra:       "二进制大对象，最大长度65535字节",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"mediumblob": {
		Name:        "mediumblob",
		FullName:    "mediumblob",
		Extra:       "二进制大对象，最大长度16777215字节",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"longblob": {
		Name:        "longblob",
		FullName:    "longblob",
		Extra:       "二进制大对象，最大长度4294967295字节",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"enum": {
		Name:        "enum",
		FullName:    "enum",
		Extra:       "枚举类型，允许定义一组固定的可能值",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"set": {
		Name:        "set",
		FullName:    "set",
		Extra:       "集合类型，允许定义一组固定的可能值，并且可以包含多个值",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
	"bit": {
		Name:        "bit",
		FullName:    "bit",
		Extra:       "位字段，存储位值，长度为M",
		DstDb:       "mysql",
		DefaultLen:  1,
		DecimalCols: 0,
	},
	"json": {
		Name:        "json",
		FullName:    "json",
		Extra:       "存储JSON数据",
		DstDb:       "mysql",
		DefaultLen:  0,
		DecimalCols: 0,
	},
}
