package ftp

import "github.com/jlaffaye/ftp"

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
