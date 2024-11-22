package httpServer

import (
	"github.com/helays/utils/ulogs"
	"github.com/helays/utils/userDb"
	"gorm.io/gorm"
	"net/http"
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
//
// User helay
// Date: 2024/11/23 0:44
//

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

// RespListsPkId 通用查询列表 主键 id
func RespListsPkId[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager ...Pager) {
	var (
		pageField, pageSizeField string
		pageSize                 int
		order                    = "id desc"
	)
	if len(pager) > 0 {
		_pager := pager[0]
		pageField = _pager.PageField
		pageSizeField = _pager.PageSizeField
		pageSize = _pager.PageSize
		order = _pager.Order
	}
	respLists(w, r, tx, respData, Pager{
		PageSize:      pageSize,
		PageSizeField: pageSizeField,
		PageField:     pageField,
		Order:         order,
	})
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
		pageField = _pager.PageField
		pageSizeField = _pager.PageSizeField
		pageSize = _pager.PageSize
		order = _pager.Order
	}
	respLists(w, r, tx, respData, Pager{
		PageSize:      pageSize,
		PageSizeField: pageSizeField,
		PageField:     pageField,
		Order:         order,
	})
}

// respLists 通用查询列表
func respLists[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager Pager) {
	var totals int64
	tx.Scopes(userDb.QueryDateTimeRange(r))
	tx.Count(&totals)
	tx.Order(pager.Order)
	if err := tx.Scopes(userDb.Paginate(r, pager.PageField, pager.PageSizeField, pager.PageSize)).Find(&respData).Error; err != nil {
		SetReturn(w, 1, "数据查询失败")
		ulogs.Error(err, r.URL.Path, r.URL.RawQuery, "respLists", "tx.Find")
		return
	}
	SetReturnData(w, 0, "成功", RespDataStruct[T]{Lists: respData, Total: totals})
}
