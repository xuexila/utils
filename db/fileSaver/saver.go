package fileSaver

import (
	"fmt"
	"github.com/helays/utils/db/fileSaver/ftp"
	"github.com/helays/utils/db/fileSaver/local"
	"github.com/helays/utils/db/fileSaver/sftp"
	"io"
	"path"
	"strings"
)

type Saver struct {
	StorageType string                                  `json:"storage_type" yaml:"storage_type" ini:"storage_type"` // 存储类型 local、sftp、ftp、hdfs等
	Root        string                                  `json:"root" yaml:"root" ini:"root"`                         // 文件跟路径，如果/开头，最终路径为/root/path，如果没有/,最终路径是current_path/root/path
	local.Local `json:"local" yaml:"local" ini:"local"` // 本地文件系统
	SftpConfig  sftp.Config                             `json:"sftp_config" yaml:"sftp_config" ini:"sftp_config"` // sftp客户端配置
	FtpConfig   ftp.Config                              `json:"ftp_config" yaml:"ftp_config" ini:"ftp_config"`    // ftp客户端配置
}

func (this Saver) Write(p string, src io.Reader, existIgnores ...bool) error {
	p = path.Join(this.Root, p)
	switch strings.ToLower(this.StorageType) {
	case "local": // 本地文件系统
		return this.Local.Write(p, src, existIgnores...)
	case "ftp": // ftp
		return this.FtpConfig.Write(p, src, existIgnores...)
	case "sftp": // sftp
		return this.SftpConfig.Write(p, src, existIgnores...)
	case "hdfs": // hdfs
	default:
		return fmt.Errorf("不支持的存储类型: %s", this.StorageType)
	}
	return nil
}
