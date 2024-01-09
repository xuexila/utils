package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"math"
	"math/rand"
	url2 "net/url"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func JsonEncode(j any) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(j)
	if err != nil {
		return nil, err
	}

	return bf.Bytes(), nil
}

// SignalHandle 系统信号
// @var funds 结束服务前，需要执行的操作
func SignalHandle(funds ...func()) {
	exitsin := make(chan os.Signal)
	signal.Notify(exitsin, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM) // 注意，syscall.SIGKILL 不能被捕获
	Log("退出信号", <-exitsin)                                                                   // 日志记录
	for _, f := range funds {
		f()
	}
	Log("各个组件关闭完成，系统即将自动关闭", os.Getpid())
	os.Exit(0)
}

// DeleteStrarr 删除字符串切片的某一个元素
func DeleteStrarr(arr []string, val string) []string {
	for index, _id := range arr {
		if _id == val {
			arr = append(arr[:index], arr[index+1:]...)
			break
		}
	}
	return arr
}

func NewId() string {
	return bson.NewObjectId().Hex()
}

// Md5 md5 函数
func Md5(s []byte) string {
	h := md5.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil))
}

// Md5string 给字符串Md5
func Md5string(s string) string {
	return Md5([]byte(s))
}

// Md5file 计算文件的Md5
func Md5file(path string) string {
	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return Md5(byt)
}

// FilePutContents 快速简易写文件
func FilePutContents(path, content string) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	CloseFile(file)
	return err
}

func FilePutContentsbytes(path string, content []byte) error {
	_path := filepath.Dir(path)
	if _, err := os.Stat(_path); err != nil {
		if err := Mkdir(_path); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	CloseFile(file)
	return err
}

// 快速简易写文件（追加）
func FileAppendContents(path, content string) error {
	_path := filepath.Dir(path)
	if _, err := os.Stat(_path); err != nil {
		if err := Mkdir(_path); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	CloseFile(file)
	return err
}

// 快速简易读取文件
func FileGetContents(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer CloseFile(file)
	return ioutil.ReadAll(file)
}

// 判断目录是否存在，否则创建目录
func Mkdir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// 将 query部分进行 url encode
func UrlEncode(url string) string {
	u, err := url2.Parse(url)
	if err != nil {
		return "-"
	}
	u.RawQuery = url2.PathEscape(u.RawQuery)
	return u.String()
}

func Int32tostring(i int32) string {
	return strconv.Itoa(int(i))
}

func Int32tobooltoint(i int32) int {
	if i > 0 {
		return 1
	}
	return 0
}

func Int64tostring(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float32tostring(f float32) string {
	f64 := float64(f)
	if math.IsNaN(f64) || math.IsInf(f64, 0) {
		return "0"
	}

	return strconv.FormatFloat(f64, 'f', 6, 64)
}

func Float64tostring(f float64) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "0"
	}
	return strconv.FormatFloat(f, 'f', 6, 64)
}

func MaxInt32(d1, d2 int32) int32 {
	if d1 > d2 {
		return d1
	}
	return d2
}

// 计算平均数
func AvgInt32(d1, d2 int32, isf bool) int32 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
}

func MinInt32(d1, d2 int32) int32 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func MaxInt64(d1, d2 int64) int64 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinInt64(d1, d2 int64) int64 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func AvgInt64(d1, d2 int64, isf bool) int64 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
}

func MaxUint64(d1, d2 uint64) uint64 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinUint64(d1, d2 uint64) uint64 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func AvgUint64(d1, d2 uint64, isf bool) uint64 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
}

func MaxFloat32(d1, d2 float32) float32 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinFloat32(d1, d2 float32) float32 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func AvgFloat32(d1, d2 float32, isf bool) float32 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
}

// 字符串转 float 64
func StrToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func Bool1time(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Booltostring 布尔转 1 0
func Booltostring(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func Uint64tostring(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func Uint16ToBytes(n int) ([]byte, error) {
	tmp := uint16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	err = binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes(), err
}

func Uint32ToBytes(n int) ([]byte, error) {
	tmp := uint32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	err = binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes(), err
}

// // 字节转换成整形
// func BytesToInt(b []byte) (int, error) {
// 	bytesBuffer := bytes.NewBuffer(b)
//
// 	var x int32
// 	err = binary.Read(bytesBuffer, binary.BigEndian, &x)
//
// 	return int(x), err
// }

// 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func BytesToUint16(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint16
	_ = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return uint16(tmp)
}

func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.LittleEndian, data); err != nil {
		return nil, nil
	}

	r, err := gzip.NewReader(b)
	if err != nil {
		// fmt.Println("[ParseGzip] NewReader error: %v, maybe data is ungzip", err)
		return data, nil
	} else {
		defer func() {
			_ = r.Close()
		}()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return undatas, nil
	}
}

// 空字符串转为 -
func EmptyString2(s string) string {
	if s = strings.TrimSpace(s); s == "" {
		return "-"
	}
	return s
}

func NumberEmptyString(s string) string {
	if s = strings.TrimSpace(s); s == "" {
		return "0"
	}
	return s
}

// Searchslice 在切片中判断某个值是否存在
func Searchslice(s string, o []string) bool {
	s = strings.TrimSpace(s)
	for _, i := range o {
		i = strings.TrimSpace(i)
		if i == s {
			return true
		}
	}
	return false
}

func SearchIntSlice(s int, arr []int) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}
func SearchInt64Slice(s int64, arr []int64) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

// Fileabs 生成文件的绝对路径
func Fileabs(cpath string) string {
	if filepath.IsAbs(cpath) {
		return cpath
	}
	return filepath.Join(Appath, cpath)
}

// MachineCode 利用硬件信息 生成token
func MachineCode() string {
	_os := strings.ToLower(runtime.GOOS)
	_arch := strings.ToLower(runtime.GOARCH)
	Log("系统版本", _os, "平台", _arch, "读取设备信息...")
	var (
		info        string
		_macode     string
		machineCode string
	)
	switch _os {
	case "linux":
		info, err = ExecShell("dmidecode")
		if err != nil {
			_macode, err = ExecShell("/bin/bash", "-c", `"/sbin/ip link" | grep link | /usr/bin/sort | /usr/bin/uniq | /usr/bin/sha256sum`)
			if err != nil {
				Error("机器码生成失败", err)
				os.Exit(1)
			}
		}
	case "darwin":
		var tmp []byte
		tmp, err = FileGetContents("/Users/helay/go/src/company/vis-device/startUp/vis.agent/run/info")
		if err != nil {
			Error("机器码生成失败", err)
			os.Exit(1)
		}
		info = string(tmp)
	default:
		os.Exit(1)
	}
	if _macode != "" {
		machineCode = Md5string(_macode)
	} else {
		cpuPreg := regexp.MustCompile(`Processor[\s\S]+?ID.+?((?:[A-Z0-9]{2} ?){8})`)

		tmp := cpuPreg.FindStringSubmatch(info)
		if len(tmp) != 2 {
			Error("系统信息获取失败")
			os.Exit(0)
		}
		cpuid := strings.TrimSpace(tmp[1])

		boardPreg := regexp.MustCompile(`Base Board[\s\S]+?Serial Number.+?([A-Z0-9]+)`)

		tmp = boardPreg.FindStringSubmatch(info)
		if len(tmp) != 2 {
			Error("系统信息获取失败")
			os.Exit(0)
		}
		boardid := tmp[1]
		// 生成机器码
		machineCode = Md5string(cpuid + Salt + boardid)
	}
	Debug("机器码", machineCode)
	return machineCode
}

// RandomString 伪随机字符串
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune
	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rd.Intn(len(letters))]
	}
	return string(b)
}
