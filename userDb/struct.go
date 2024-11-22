package userDb

import "github.com/helays/utils/dataType"

type Dbbase struct {
	DbType       string               `ini:"db_type" yaml:"db_type" json:"db_type" gorm:"type:varchar(32);not null;comment:数据库类型，mysql|pg等"` // 数据库类型 mysql/pg
	Host         dataType.StringArray `ini:"host" yaml:"host" json:"host" gorm:"not null;comment:连接信息"`
	User         string               `ini:"user" yaml:"user" json:"user" gorm:"type:varchar(256);not null;comment:数据库用户"`
	Pwd          string               `ini:"pwd" yaml:"pwd" json:"pwd" gorm:"type:varchar(512);not null;comment:数据库密码"`
	Encrypt      string               `ini:"encrypt" yaml:"encrypt" json:"encrypt" gorm:"type:varchar(64);not null;default:none;comment:密码加密存储方式，默认none明文"` // 是否加密
	Dbname       string               `ini:"dbname" yaml:"dbname" json:"dbname" gorm:"type:varchar(128);not null;comment:默认连接的库"`
	MaxIdleConns int                  `ini:"max_idle_conns" yaml:"maxIdleConns" json:"maxIdleConns" gorm:"type:int;not null;default:2;comment:最大空闲连接数"`
	MaxOpenConns int                  `ini:"max_open_conns" yaml:"maxOpenConns" json:"maxOpenConns" gorm:"type:int;not null;default:10;comment:最大连接数"`
	Remark       string               `ini:"remark" yaml:"remark" json:"remark" gorm:"type:varchar(256);not null;uniqueIndex;comment:备注信息，只能唯一"`
}
