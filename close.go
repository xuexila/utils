package utils

import (
	"database/sql"
	"github.com/IBM/sarama"
	"github.com/colinmarc/hdfs/v2"
	"github.com/garyburd/redigo/redis"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	rds "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
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

// CloseFtpClient ftp连接退出关闭
func CloseFtpClient(conn *ftp.ServerConn) {
	if conn != nil {
		_ = conn.Logout()
		_ = conn.Quit()
	}
}

func CloseFtpResponse(raw *ftp.Response) {
	if raw != nil {
		_ = raw.Close()
	}
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

func CloseSsh(conn *ssh.Client) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseSftp(conn *sftp.Client) {
	if conn != nil {
		_ = conn.Close()
	}
}

// CloseRedisConn github.com/garyburd/redigo/redis
func CloseRedisConn(conn redis.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseRdsConn(conn *rds.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseKafkaPartition(partition sarama.PartitionConsumer) {
	if partition == nil {
		return
	}
	Checkerr(partition.Close(), "CloseKafkaPartition")
}

func CloseFile(file *os.File) {
	if file != nil {
		_ = file.Close()
	}
}

func CloseHdfsFile(f *hdfs.FileWriter) {
	if f != nil {
		_ = f.Close()
	}
}

func CloseMultipartWriter(w *multipart.Writer) {
	if w != nil {
		_ = w.Close()
	}
}

func CloseSftpFile(file *sftp.File) {
	if file != nil {
		_ = file.Close()
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
