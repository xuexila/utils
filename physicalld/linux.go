//go:build linux
// +build linux

package physicalld

import (
	"bytes"
	"code.helay.net/helay/utils/crypto/md5"
	"code.helay.net/helay/utils/fileTools"
	"compress/gzip"
	"encoding/binary"
	"io"
	"os"
	"os/exec"
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
// Date: 2022/10/12 0:04
//

func GenPhysicalld() (string, error) {
	return linux()
}

// 获取linux系统的机器码
func linux() (string, error) {
	dmidecodes, err := bindataRead(dmidecode, "Dmidecode")
	if err != nil {
		return "", err
	}
	if binaryWrite(dmidecodes, dmidecodePath) {
		defer os.Remove(dmidecodePath)
		if err := os.Chmod(dmidecodePath, 0700); err != nil {
			return "", err
		}

		info, err := execShell(dmidecodePath + " -t system|grep 'Serial Number'|head -1 ; " + dmidecodePath + " -t Processor|grep ID |head -1 ; " + dmidecodePath + " -t system|grep UUID")
		if err == nil && info != "" {
			return md5.Md5string(info + string(salt)), nil
		}
	}
	info, err := execShell("ip link | grep link | /usr/bin/sort | /usr/bin/uniq |/usr/bin/sha256sum")
	if err != nil {
		return "", err
	}
	return md5.Md5string(info + string(salt)), nil
}

func binaryWrite(data []byte, filename string) bool {
	if filename == "" {
		return false
	}
	file, err := os.Create(filename)
	defer fileTools.CloseFile(file)
	if err != nil {
		return false
	}
	var binBuf bytes.Buffer
	_ = binary.Write(&binBuf, binary.LittleEndian, data)
	b := binBuf.Bytes()
	_, err = file.Write(b)
	return err == nil
}

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()
	if err != nil {
		return nil, err
	}
	if clErr != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func execShell(s string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", s)
	var out = new(bytes.Buffer)
	cmd.Stdout = out
	err := cmd.Run()
	return out.String(), err
}
