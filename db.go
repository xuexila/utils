package utils

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
			return db.Where(query, "%"+field+"%")
		}
		return db.Where(query, field)
	}
}
