package hdfsClose

import "github.com/colinmarc/hdfs/v2"

func CloseHdfsFile(f *hdfs.FileWriter) {
	if f != nil {
		_ = f.Close()
	}
}

func CloseHdfs(conn *hdfs.Client) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseHdfsFileReader(file *hdfs.FileReader) {
	if file != nil {
		_ = file.Close()
	}
}
