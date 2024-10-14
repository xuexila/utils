package sftp_close

import "github.com/pkg/sftp"

func CloseSftp(conn *sftp.Client) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseSftpFile(file *sftp.File) {
	if file != nil {
		_ = file.Close()
	}
}
