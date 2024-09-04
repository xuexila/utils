package userDb

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var nullSqlConn map[string]*gorm.DB

// GetRawSql 生成sql的通用函数
func GetRawSql(f func(dTx *gorm.DB) *gorm.DB, dbTypes ...string) (string, []any) {
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

// GetRawSqlByDb 生成sql的通用函数,db为数据库连接
func GetRawSqlByDb(f func(dTx *gorm.DB) *gorm.DB, db *gorm.DB) (string, []any) {
	query := f(db).Statement
	return query.SQL.String(), query.Vars
}

func CloseDb(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseMysqlRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}

// Deprecated: As of utils v1.1.0, this value is simply [utils.CloseDb].
func CloseMysql(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseStmt(stmt *sql.Stmt) {
	if stmt != nil {
		_ = stmt.Close()
	}
}
