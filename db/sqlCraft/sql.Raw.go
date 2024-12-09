package sqlCraft

import (
	"database/sql"
	"encoding/json"
	"github.com/helays/utils/db/userDb"
	"github.com/helays/utils/http/httpExport/httpExportExcel"
	"github.com/helays/utils/http/httpServer"
	"github.com/helays/utils/ulogs"
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
// Date: 2024/11/23 16:59
//

func (this SqlFilter) RunSql(w http.ResponseWriter, r *http.Request, tx *gorm.DB) {
	if this.Sql == "" {
		httpServer.SetReturnCode(w, r, 500, "无可执行sql")
		return
	}
	var err error
	if this.Type == "exec" {
		exx := tx.Exec(this.Sql, this.Args...)
		err = exx.Error
		if err != nil {
			httpServer.SetReturnError(w, r, err, 500, "执行原生sql失败")
			return
		}
		httpServer.SetReturnCode(w, r, 0, "执行原生sql成功", exx.RowsAffected)
		return
	}
	var rows *sql.Rows
	rows, err = tx.Raw(this.Sql, this.Args...).Rows()
	defer userDb.CloseMysqlRows(rows)
	if err != nil {
		httpServer.SetReturnError(w, r, err, 500, "执行原生sql失败")
		return
	}
	this.Response(w, tx, rows)
}

func (this SqlFilter) Response(w http.ResponseWriter, uTx *gorm.DB, rows *sql.Rows) (bool, int) {
	if this.Export {
		if err := httpExportExcel.RespExcelOrCsv(w, rows, this.FileType, this.ExportHeader); err != nil {
			ulogs.Error(err, "原生SQL导出excel或csv失败")
			return false, 0
		}
		return true, 0
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"code":0,"msg":"查询成功","data":[`))
	isf := true
	for rows.Next() {
		if !isf {
			_, _ = w.Write([]byte(","))
		}
		isf = false
		var result map[string]any
		if err := uTx.ScanRows(rows, &result); err != nil {
			ulogs.Error(err, "原生sql查询 scan失败", this.Sql)
			continue
		}
		ulogs.Checkerr(json.NewEncoder(w).Encode(result), "原生sql查询 响应失败")
	}
	_, _ = w.Write([]byte(`]}`))
	return true, 0
}
