package httpServerWithDb

import (
	"github.com/helays/utils/db/userDb"
	"github.com/helays/utils/http/httpServer"
	"github.com/helays/utils/tools"
	"github.com/helays/utils/ulogs"
	"gorm.io/gorm"
	"net/http"
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
		pageField, pageSizeField string
		pageSize                 int
		order                    = "row_id desc"
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
		pageField, pageSizeField string
		pageSize                 int
		order                    = "id desc"
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
