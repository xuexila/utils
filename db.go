package utils

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"reflect"
	"strconv"
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
// User helei
// Date: 2023/9/1 11:19
//

// Paginate 分页通用部分
func Paginate(r *http.Request, pageField, pageSizeField string, pageSize int) func(db *gorm.DB) *gorm.DB {
	if pageField == "" {
		pageField = "pageNo"
	}
	if pageSizeField == "" {
		pageSizeField = "pageSize"
	}
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(r.URL.Query().Get(pageField))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get(pageSizeField))
		if limit < 1 {
			limit = pageSize
		}
		tx := db
		if r.URL.Query().Get("rall") != "1" {
			tx.Offset((page - 1) * limit).Limit(limit)
		}
		_sort := r.URL.Query().Get("sort")
		if _sort != "" && !specialChartPreg.MatchString(_sort) {
			if _sort[0] == '-' {
				tx.Order(_sort[1:] + " desc")
			} else {
				tx.Order(_sort)
			}
		}
		return tx
	}
}

// FilterWhereString 过滤string 条件
func FilterWhereString(r *http.Request, query string, field string, like bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		value := r.URL.Query().Get(field)
		if value == "" {
			return db
		}
		if like {
			return db.Where(query, "%"+value+"%")
		}
		return db.Where(query, value)
	}
}

// FilterWhereStruct 通过结构体 自动映射查询字段
func FilterWhereStruct(s any, alias string, enableDefault bool, r *http.Request, likes ...map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		t := reflect.TypeOf(s)
		if t.Kind() != reflect.Struct {
			return db
		}
		query := r.URL.Query()
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Type.String() != "int" && t.Field(i).Type.String() != "string" {
				continue
			}
			tagName := t.Field(i).Tag.Get("json")
			if tagName == "" {
				continue
			}
			val := query.Get(strings.Split(tagName, ",")[0])
			if val == "" {
				if !enableDefault {
					continue
				}
				// 如果没有传值，判断是否有默认值
				val = t.Field(i).Tag.Get("default")
				if val == "" {
					continue
				}
			}

			tableName := alias
			if tableName == "" {
				tableName = SnakeString(t.Name())
			}
			// 这里还需要解析出字段本身的名字，去数据库进行查询，通过将结构体转成蛇形方式。
			fieldName := tableName + "." + SnakeString(t.Field(i).Name)
			if t.Field(i).Type.String() == "int" {
				db.Where(fieldName+"=?", val)
			} else {
				lastVal := val
				if t.Field(i).Tag.Get("dblike") == "%" {
					lastVal = "%" + val + "%"
				}
				if len(likes) > 0 {
					if custom, ok := likes[0][fieldName]; ok {
						switch custom {
						case "%%":
							lastVal = "%" + val + "%"
						case "-%":
							lastVal = val + "%"
						case "%-":
							lastVal = "%" + val
						default:
							lastVal = val
						}
					}
				}
				db.Where(fieldName+" like ?", lastVal)
			}
		}
		return db
	}

}

// QueryDateTimeRange 时间区间查询
func QueryDateTimeRange(r *http.Request, filed ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sTime := r.URL.Query().Get("begin_time")
		eTime := r.URL.Query().Get("end_time")
		sField := "create_time"
		if len(filed) > 0 {
			sField = filed[0]
		}
		if sTime != "" {
			db.Where(sField+" > ?", sTime)
		}
		if eTime != "" {
			db.Where(sField+" <= ?", eTime)
		}
		return db
	}
}

// AutoCreateTableWithStruct 根据结构体判断是否需要创建表
func AutoCreateTableWithStruct(db *gorm.DB, tb any, errmsg string) {
	t := reflect.TypeOf(tb)
	if t.Kind() != reflect.Struct {
		return
	}
	if !db.Migrator().HasTable(tb) {
		DieCheckerr(db.Debug().AutoMigrate(tb), errmsg)
	}
	// 如果表存在，在判断结构体中是否有新增字段，如果有，就自动改变表
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("gorm")
		if tag == "" {
			continue
		}
		if tag == "-:all" || tag == "-:migration" {
			continue
		}
		column := SnakeString(t.Field(i).Name)
		for _, item := range strings.Split(tag, ";") {
			if !strings.HasPrefix(item, "column") {
				continue
			}
			column = item[7:]
		}
		if !db.Migrator().HasColumn(tb, column) {
			DieCheckerr(db.Debug().AutoMigrate(tb), errmsg)
			break // 创建一次就行了
		}
	}
}

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
		case "sqlite":
			dialector = sqlite.Open("")
		case "sqlserver":
			dialector = sqlserver.Open("")
		case "tidb":
			dialector = mysql.Open("")
		case "clickhouse":
			dialector = clickhouse.Open("")

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
