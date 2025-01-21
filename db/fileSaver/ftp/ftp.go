package ftp

import (
	"fmt"
	"github.com/helays/utils/close/ftpClose"
	"github.com/jlaffaye/ftp"
	"io"
	"path"
)

// Config sftp 配置
type Config struct {
	Host string `json:"host" yaml:"host" ini:"host"` // ftp地址:端口
	User string `json:"user" yaml:"user" ini:"user"`
	Pwd  string `json:"pwd" yaml:"pwd" ini:"pwd"` // 密码
	// 这部分是ftp的
	Epsv int `ini:"epsv" yaml:"epsv" json:"epsv,omitempty" gorm:"type:int;not null;default:0;comment:是否启用加密"` // ftp 连接模式，0 被动模式 1 主动模式
}

// Write 写入文件
func (this Config) Write(p string, src io.Reader, existIgnores ...bool) error {
	ftpClient, err := this.Login()
	if err != nil {
		return err
	}
	defer ftpClose.CloseFtpClient(ftpClient)
	filePath, err := SetPath(ftpClient, p)
	if err != nil {
		return err
	}
	// 判断是否需要覆盖写入
	if len(existIgnores) > 0 && existIgnores[0] {
		if ok, err := Exist(ftpClient, filePath); ok {
			return nil
		} else if err != nil {
			return err
		}
	}
	dir := path.Dir(filePath)
	// 首先判断这个路径是否存在，然后创建
	if err = Mkdir(ftpClient, dir); err != nil {
		return err
	}
	if err = ftpClient.Stor(filePath, src); err != nil {
		return fmt.Errorf("写入文件%s失败：%s", filePath, err.Error())
	}
	return nil
}

// Login ftp登录
func (this Config) Login() (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(this.Host, ftp.DialWithDisabledEPSV(this.Epsv == 1))
	if err != nil {
		return nil, fmt.Errorf("ftp连接失败：%s", err.Error())
	}
	if err = conn.Login(this.User, this.Pwd); err != nil {
		return nil, fmt.Errorf("ftp登录失败：%s", err.Error())
	}
	return conn, nil
}
