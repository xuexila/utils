package utils

import (
	"database/sql"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
)

func CloseResp(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	_ = resp.Close
}

func CloseReq(resp *http.Request) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
}

func CloseConn(conn net.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}
func CloseUdpConn(conn *net.UDPConn) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseFile(file *os.File) {
	if file != nil {
		_ = file.Close()
	}
}

func CloseMultipartWriter(w *multipart.Writer) {
	if w != nil {
		_ = w.Close()
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

func CloseDb(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseStmt(stmt *sql.Stmt) {
	if stmt != nil {
		_ = stmt.Close()
	}
}

func Closews(conn *websocket.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseMultipartFile(f multipart.File) {
	if f != nil {
		_ = f.Close()
	}
}

func Closeresponse(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
}
