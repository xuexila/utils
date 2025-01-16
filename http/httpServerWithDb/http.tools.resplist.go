package httpServerWithDb

import (
	"github.com/helays/utils/db/userDb"
	"github.com/helays/utils/http/httpServer"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/tools"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type Pager struct {
	PageSize      int    `ini:"page_size" yaml:"page_size" json:"page_size"` // 系统默认查询数量
	PageSizeField string `ini:"page_size_field" yaml:"page_size_field" json:"page_size_field"`
	PageField     string `ini:"page_field" yaml:"page_field" json:"page_field"`
	Order         string
}

type RespDataStruct[T any] struct {
	Lists T     `json:"lists"`
	Total int64 `json:"total"`
}

// respLists 通用查询列表
func respLists[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager Pager) {
	var totals int64
	tx.Scopes(userDb.QueryDateTimeRange(r))
	tx.Count(&totals)
	tx.Order(pager.Order)
	if err := tx.Scopes(userDb.Paginate(r, pager.PageField, pager.PageSizeField, pager.PageSize)).Find(&respData).Error; err != nil {
		httpServer.SetReturn(w, 1, "数据查询失败")
		ulogs.Error(err, r.URL.Path, r.URL.RawQuery, "respLists", "tx.Find")
		return
	}
	httpServer.SetReturnData(w, 0, "成功", RespDataStruct[T]{Lists: respData, Total: totals})
}

// RespListsPkRowId 通用查询列表 主键 row_id
func RespListsPkRowId[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager ...Pager) {
	var (
		pageField     = "p"  // 页面默认字段
		pageSizeField = "pn" // 页面呈现数量默认字段
		pageSize      = 30   // 每页默认数量
		order         = "row_id desc"
	)
	if len(pager) > 0 {
		_pager := pager[0]
		pageField = tools.Ternary(_pager.PageField == "", pageField, _pager.PageField)
		pageSizeField = tools.Ternary(_pager.PageSizeField == "", pageSizeField, _pager.PageSizeField)
		pageSize = tools.Ternary(_pager.PageSize < 1, pageSize, _pager.PageSize)
		order = tools.Ternary(_pager.Order == "", order, _pager.Order)
	}
	respLists(w, r, tx, respData, Pager{
		PageSize:      pageSize,
		PageSizeField: pageSizeField,
		PageField:     pageField,
		Order:         order,
	})
}

// RespListsPkId 通用查询列表 主键 id
func RespListsPkId[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager ...Pager) {
	var (
		pageField     = "p"
		pageSizeField = "pn"
		pageSize      = 30
		order         = "id desc"
	)
	if len(pager) > 0 {
		_pager := pager[0]
		pageField = tools.Ternary(_pager.PageField == "", pageField, _pager.PageField)
		pageSizeField = tools.Ternary(_pager.PageSizeField == "", pageSizeField, _pager.PageSizeField)
		pageSize = tools.Ternary(_pager.PageSize < 1, pageSize, _pager.PageSize)
		order = tools.Ternary(_pager.Order == "", order, _pager.Order)
	}
	respLists(w, r, tx, respData, Pager{
		PageSize:      pageSize,
		PageSizeField: pageSizeField,
		PageField:     pageField,
		Order:         order,
	})
}

// ListMethodGet 通用查询函数（使用泛型），请求方式是Get
func ListMethodGet[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, c userDb.QueryConfig, p Pager) {
	var (
		list []T
		mod  T
	)
	_tx := tx.Debug().Model(list)
	_tx.Scopes(userDb.FilterWhereStruct(mod, "", false, r))
	if c.SelectQuery != nil {
		_tx.Select(c.SelectQuery, c.SelectArgs...)
	}
	if c.Query != nil {
		_tx.Where(c.Query, c.Args...)
	}
	switch strings.ToLower(c.Pk) {
	case "id":
		RespListsPkId(w, r, _tx, list, p)
	case "row_id":
		RespListsPkRowId(w, r, _tx, list, p)
	default:
		return
	}

}
