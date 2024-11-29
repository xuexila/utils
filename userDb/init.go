package userDb

import (
	"errors"
	"fmt"
	"github.com/helays/utils/tools"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/url"
	"strings"
	"time"
)

func (this Dbbase) Dsn() string {
	dsn := url.URL{
		User: url.UserPassword(this.User, this.Pwd),
		Host: strings.Join(this.Host, ","),
		Path: this.Dbname,
	}
	query := dsn.Query()
	switch this.DbType {
	case "pg":
		dsn.Scheme = "postgres"
		query.Set("TimeZone", "Asia/Shanghai")
		if this.Schema != "" {
			query.Set("search_path", this.Schema)
		}
	case "mysql":
		//dsn.Scheme = "mysql" // mysql 不需要这个
		dsn.Host = fmt.Sprintf("tcp(%s)", dsn.Host)
		query.Set("charset", "utf8mb4")
		query.Set("parseTime", "True")
		query.Set("loc", "Local")
	}
	dsn.RawQuery = query.Encode()
	return dsn.String()
}

// InitDb 连接数据库
func InitDb(c Dbbase) (*gorm.DB, error) {
	var (
		dsn       = c.Dsn()
		dialector gorm.Dialector
	)
	switch c.DbType {
	case "pg":
		//postgres://user:password@host1:port1,host2:port2/database?target_session_attrs=read-write&TimeZone=Asia/Shanghai
		//dsn = "postgres://" + c.User + ":" + c.Pwd + "@" + strings.Join(c.Host, ",") + "/" + c.Dbname + "?TimeZone=Asia/Shanghai"
		dialector = postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		})
	case "mysql":
		dsn = strings.TrimLeft(dsn, "//")
		//fmt.Println(dsn)
		//dsn = c.User + ":" + c.Pwd + "@tcp(" + strings.Join(c.Host, ",") + ")/" + c.Dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
		//fmt.Println(dsn)
		dialector = mysql.New(mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		})
	default:
		return nil, errors.New("不支持的数据库")
	}
	_db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		return nil, err
	}
	_sqlDb, err := _db.DB()
	if err != nil {
		CloseGormDb(_db)
		return nil, err
	}
	_sqlDb.SetMaxIdleConns(tools.Ternary(c.MaxIdleConns < 1, 2, c.MaxIdleConns)) // 设置连接池中空闲连接的最大数量
	_sqlDb.SetMaxOpenConns(tools.Ternary(c.MaxOpenConns < 1, 5, c.MaxOpenConns)) // 设置打开数据库连接的最大数量
	_sqlDb.SetConnMaxIdleTime(time.Hour)                                         // 连接空闲1小时候将失效
	return _db, nil
}

func CloseGormDb(_db *gorm.DB) {
	if _db == nil {
		return
	}
	if sqlDb, err := _db.DB(); err == nil {
		CloseDb(sqlDb)
	}
}
