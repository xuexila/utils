package httpServerWithDb

import (
	"github.com/helays/utils/config"
	"github.com/helays/utils/db/userDb"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// ListMethodGet 是一个通用的列表查询方法，用于根据不同的条件获取数据库中的记录。
// 它使用了泛型 T，允许任何类型的列表查询。
// 参数:
//
//	w http.ResponseWriter: 用于写入HTTP响应。
//	r *http.Request: 包含当前HTTP请求的详细信息。
//	tx *gorm.DB: GORM数据库连接对象，用于执行数据库操作。
//	c userDb.QueryConfig: 查询配置，包含了查询所需的配置信息，如选择查询、条件查询等。
//	p Pager: 分页配置，用于指定查询的分页信息。
func ListMethodGet[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, c userDb.QueryConfig, p Pager) {
	if config.Dbg {
		tx = tx.Debug()
	}
	_tx := tx.Model(new(T))
	_tx.Scopes(userDb.FilterWhereStruct(new(T), "", false, r))
	if c.SelectQuery != nil {
		_tx.Select(c.SelectQuery, c.SelectArgs...)
	}
	if c.Query != nil {
		_tx.Where(c.Query, c.Args...)
	}
	var list []T
	switch strings.ToLower(c.Pk) {
	case "id":
		RespListsPkId(w, r, _tx, list, p)
	case "row_id":
		RespListsPkRowId(w, r, _tx, list, p)
	default:
		return
	}
}
