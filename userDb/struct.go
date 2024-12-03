package userDb

import "github.com/helays/utils/dataType"

var (
	// SupportedDbType 支持的数据库类型
	SupportedDbType = []map[string]string{
		{"type": "mysql", "value": "mysql"},
		{"type": "pg", "value": "postgres"},
		{"type": "kafka", "value": "kafka"},
		{"type": "ftp", "value": "ftp"},
		{"type": "sftp", "value": "sftp"},
	}
	// FTPEpsv ftp模式
	FTPEpsv = []map[string]any{
		{"type": 0, "value": "被动模式"},
		{"type": 1, "value": "主动模式"},
		{"type": 2, "value": "自动模式"},
	}
	// Authentication 认证方式
	Authentication = []map[string]string{
		{"type": "password", "value": "密码"},
		{"type": "public_key", "value": "密钥"},
	}
)

type Dbbase struct {
	Remark string `ini:"remark" yaml:"remark" json:"remark" gorm:"type:varchar(256);not null;uniqueIndex;comment:备注信息，只能唯一"`
	DbType string `ini:"db_type" yaml:"db_type" json:"db_type" gorm:"type:varchar(32);not null;index;comment:数据库类型，mysql|pg|kafka等"` // 数据库类型 mysql/pg

	// 这部分是公用的
	Host    dataType.StringArray `ini:"host" yaml:"host" json:"host" gorm:"not null;comment:连接信息"`
	User    string               `ini:"user" yaml:"user" json:"user" gorm:"type:varchar(256);not null;default:'';comment:数据库用户"`
	Pwd     string               `ini:"pwd" yaml:"pwd" json:"pwd" gorm:"type:text;comment:数据库密码"`
	Encrypt string               `ini:"encrypt" yaml:"encrypt" json:"encrypt" gorm:"type:varchar(64);not null;default:none;comment:密码加密存储方式，默认none明文"` // 是否加密

	// 这部分是数据库独有
	Dbname       string `ini:"dbname" yaml:"dbname" json:"dbname" gorm:"type:varchar(128);not null;index;default:'';comment:默认连接的库"`
	Schema       string `ini:"schema" yaml:"schema" json:"schema" gorm:"type:varchar(128);not null;default:'';comment:数据库模式"`
	MaxIdleConns int    `ini:"max_idle_conns" yaml:"max_idle_conns" json:"max_idle_conns" gorm:"type:int;not null;default:2;comment:最大空闲连接数"`
	MaxOpenConns int    `ini:"max_open_conns" yaml:"max_open_conns" json:"max_open_conns" gorm:"type:int;not null;default:10;comment:最大连接数"`

	// 这部分是kafka独有
	MsgRole   string `ini:"msg_role" yaml:"msg_role" json:"msg_role" gorm:"type:varchar(32);not null;index;default:'';comment:消息角色，syncProducer|asyncProducer|consumer"`
	Version   string `ini:"version" yaml:"version" json:"version" gorm:"type:varchar(32);not null;default:'';comment:kafka版本"`
	Sasl      int    `ini:"sasl" yaml:"sasl" json:"sasl" gorm:"type:int;not null;default:0;comment:是否启用加密"`
	Mechanism string `ini:"mechanism" yaml:"mechanism" json:"mechanism" gorm:"type:varchar(32);not null;default:'';comment:加密方式"`

	// 这部分是ftp的
	Epsv int `ini:"epsv" yaml:"epsv" json:"epsv" gorm:"type:int;not null;default:0;comment:是否启用加密"` // ftp 连接模式，0 被动模式 1 主动模式 2 自动
	// 这部分是sftp的，
	Authentication string `ini:"authentication" yaml:"authentication" json:"authentication" gorm:"type:varchar(32);not null;default:'';comment:认证方式"` // 密码或者 私钥认证
}
