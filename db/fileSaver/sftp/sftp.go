package sftp

import (
	"fmt"
	"github.com/helays/utils/close/vclose"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"path"
)

// Config sftp 配置
type Config struct {
	Host           string `json:"host" yaml:"host" ini:"host"` // 路径
	User           string `json:"user" yaml:"user" ini:"user"`
	Pwd            string `json:"pwd" yaml:"pwd" ini:"pwd"`                                  // 密码|密钥
	Authentication string `json:"authentication" yaml:"authentication" ini:"authentication"` // 认证方式 ，默认passwd,可选public_key
}

// Write 写入文件
func (this Config) Write(p string, src io.Reader, existIgnores ...bool) error {
	sshClient, sftpClient, err := this.Login()
	defer func() {
		vclose.Close(sftpClient)
		vclose.Close(sshClient)
	}()
	if err != nil {
		return err
	}
	filePath, err := SetPath(sftpClient, p)
	if err != nil {
		return err
	}
	// 判断是否需要覆盖写入
	if len(existIgnores) > 0 && existIgnores[0] {
		if ok, err := Exist(sftpClient, filePath); ok {
			return nil
		} else if err != nil {
			return err
		}
	}

	dir := path.Dir(filePath)
	// 首先判断这个路径是否存在，然后创建
	if err = Mkdir(sftpClient, dir); err != nil {
		return err
	}
	// 文件夹存在后，就开始创建文件
	file, err := sftpClient.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件%s失败：%s", filePath, err.Error())
	}
	defer vclose.Close(file)
	if _, err = io.Copy(file, src); err != nil {
		return fmt.Errorf("写入文件%s失败：%s", filePath, err.Error())
	}
	return nil
}

// Login ssh登录
// @return *ssh.Client, *sftp.Client, error
func (this Config) Login() (*ssh.Client, *sftp.Client, error) {
	// 首先连接 ssh client
	clientConfig := &ssh.ClientConfig{
		User:            this.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var auth ssh.AuthMethod
	if this.Authentication == "password" {
		auth = ssh.Password(this.Pwd)
	} else {
		var signer ssh.Signer
		signer, err := ssh.ParsePrivateKey([]byte(this.Pwd))
		if err != nil {
			return nil, nil, fmt.Errorf("ssh密钥解析失败：%s", err.Error())
		}
		auth = ssh.PublicKeys(signer)
	}
	clientConfig.Auth = []ssh.AuthMethod{auth}
	sshClient, err := ssh.Dial("tcp", this.Host, clientConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("ssh连接失败：%s", err.Error())
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, nil, fmt.Errorf("sftp连接失败：%s", err.Error())
	}
	return sshClient, sftpClient, nil
}
