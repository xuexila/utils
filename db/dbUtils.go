package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var nullSqlConn map[string]*gorm.DB

// GetRawSql 生成sql的通用函数
func GetRawSql(f func(tx *gorm.DB) *gorm.DB, dbTypes ...string) (string, []any) {
	dbType := "pg"
	if len(dbTypes) > 1 {
		dbType = dbTypes[0]
	}
	db, ok := nullSqlConn[dbType]
	if !ok {
		var dialector gorm.Dialector
		switch dbType {
		case "pg":
			dialector = postgres.Open("")
		case "mysql":
			dialector = mysql.Open("")
			//case "sqlite":
			//	dialector = sqlite.Open("")
			//case "sqlserver":
			//	dialector = sqlserver.Open("")
			//case "tidb":
			//	dialector = mysql.Open("")
			//case "clickhouse":
			//	dialector = clickhouse.Open("")

		}
		db, _ = gorm.Open(dialector, &gorm.Config{
			DryRun:                                   true,
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
			DisableAutomaticPing:                     true,
		})
	}

	query := f(db).Statement
	return query.SQL.String(), query.Vars
}
