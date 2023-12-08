package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/colinmarc/hdfs"
	"github.com/garyburd/redigo/redis"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	rds "github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	url2 "net/url"
	"os"
	"os/exec"
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
func SignalHandle() {
	exitsin := make(chan os.Signal)
	signal.Notify(exitsin, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM) // 注意，syscall.SIGKILL 不能被捕获
	Log("退出信号", <-exitsin)                                                                   // 日志记录
	os.Exit(0)
}

// Play 公共函数文件
func Play(path string, w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()
	if err != nil {
		Error("视频不存在", path)
		http.NotFound(w, r)
		return
	}

	ranges := int64(0)
	rangeEnd := int(0)
	rangeSwap := strings.TrimSpace(r.Header.Get("Range"))
	if rangeSwap != "" {
		rangeSwap = rangeSwap[6:]
		rangListSwap := strings.Split(rangeSwap, "-")
		if len(rangeSwap) >= 1 {
			if num, err := strconv.Atoi(rangListSwap[0]); err == nil {
				ranges = int64(num)
			}
			if len(rangeSwap) > 1 {
				if num, err := strconv.Atoi(rangListSwap[1]); err == nil {
					rangeEnd = int(num)
				}
			}
		}
	}

	var (
		fileSize int
		//tmpF     []byte
	)
	fType := strings.ToLower(filepath.Ext(path)[1:])
	fInfo, err := f.Stat()
	if err != nil {
		Forbidden(w, "403 Forbidden!")
		return
	}

	//if fType=="mp4" {
	//	tmpF = GetMP4Duration(f)
	//	fileSize= len(tmpF)
	//	if rangeSwap!="" && ranges>0 {
	//		tmpF=tmpF[ranges:]
	//	}
	//}else{
	//
	//	if rangeSwap!="" && ranges>0 {
	//		_, _ = f.Seek(ranges, 0)
	//	}
	//	fileSize=int(fInfo.Size())
	//}
	//GetMP4Duration(f)
	if rangeSwap != "" && ranges > 0 {
		_, _ = f.Seek(ranges, 0)
	}
	fileSize = int(fInfo.Size())

	totalSize := fileSize
	if rangeSwap != "" && rangeEnd > 0 {
		totalSize = rangeEnd
	}
	total := strconv.Itoa(fileSize)

	mime := MimeMap[fType]
	if mime == "" {
		mime = "text/html;charset=utf-8"
	}
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.Itoa(totalSize-int(ranges)))
	w.Header().Set("Last-Modified", fInfo.ModTime().Format(time.RFC822))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Connection", "close")
	w.Header().Set("Etag", `W/"`+strconv.FormatInt(fInfo.ModTime().Unix(), 16)+`-`+strconv.FormatInt(fInfo.Size(), 16)+`"`)
	if rangeSwap != "" {
		w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(ranges))+"-"+strconv.Itoa(totalSize-1)+"/"+total)
		w.WriteHeader(206)
	} else {
		w.WriteHeader(200)
	}
	if fType == "mp4" {
		_byt, _ := ioutil.ReadAll(f)
		_, _ = w.Write(_byt)
		return
	}
	_, _ = io.Copy(w, f)
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

// md5 函数
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

func ErrorReturn(i ...interface{}) bool {
	Error(i...)
	return false
}

// Error 打印错误信息
func Error(i ...interface{}) {
	log.SetPrefix("")
	log.SetOutput(os.Stderr)
	//_lst:=i[len(i)-1]
	//fmt.Println("ss",_lst==nil)
	log.Println(i...)
}

func ReqError(r *http.Request, i ...any) {
	log.SetPrefix("[" + Getip(r) + "] ")
	log.SetOutput(os.Stderr)
	log.Println(i...)
}

// Recover 捕获系统异常
func Recover() {
	if r := recover(); r != nil {
		Error("系统异常，捕获结果", r)
	}
}

// Log 打印正确日志。
func Log(i ...interface{}) {
	// log.SetPrefix("[用户日志]")
	log.SetOutput(os.Stdout)
	log.Println(i...)
}

func Debug(i ...interface{}) {
	if Dbg {
		Log("[debug]", i)
	}
}

// Checkerr 检查错误
func Checkerr(err error, i ...interface{}) {
	if err == nil {
		return
	}
	Error(i, err)
}

// DieCheckerr 检查错误，打印并输出错误信息
func DieCheckerr(err error, i ...interface{}) {
	if err == nil {
		return
	}
	Error(i, err)
	os.Exit(1)
}

// ReturnCheckerr 检查错误，有异常就返回false
func ReturnCheckerr(err error, i ...interface{}) bool {
	if err == nil {
		return true
	}
	Error(i, err)
	return false
}

func Pfunc(a ...interface{}) {
	// log.SetPrefix("[用户异常]")
	log.SetOutput(os.Stdout)
	log.Println(a...)
}

// ExecShell 执行 shell语句
func ExecShell(name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = bytes.NewBuffer(nil)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	Log("shell 执行命令", cmd.String())

	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}

func ExecShellQuit(name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = bytes.NewBuffer(nil)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}

// 可手动结束的 shell命令
func ExecCtlShell(stop chan byte, name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = new(bytes.Buffer)
		end = make(chan byte)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	Log("shell 执行命令", cmd.String())
	go func(cmd *exec.Cmd, stop, end chan byte) {
		for {
			select {
			case <-stop:
				if cmd == nil {
					Error("分析EXEC 不存在", cmd.String())
					continue
				}
				if cmd.Process == nil {
					Error("分析 Process不存在", cmd.String())
					continue
				}
				if err := cmd.Process.Kill(); err != nil {
					Error("手动结束进程失败", err)
				} else {
					Log("手动结束命令", cmd.String())
				}

			case <-end:
				Log("shell 执行完成", cmd.String())
				return
			}
		}
	}(cmd, stop, end)
	defer func() {
		go func() {
			end <- 1
		}()

	}()
	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}

// 快速简易写文件
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

// 设置系统返回403
func Forbidden(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusForbidden)
	_html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>403 Forbidden!</title></head><body><h3 style="text-align:center">` + msg + `</h3></body></html>`
	_, _ = fmt.Fprintln(w, _html)
	return
}

// 设置返回 404
func NotFound(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusNotFound)
	_html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>404 Not Found!</title></head><body><h3 style="text-align:center">` + msg + `</h3></body></html>`
	_, _ = fmt.Fprintln(w, _html)
	return
}

// 设置返回 405
func MethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
}

func Closeresponse(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
}

// 生成文件的绝对路径
func Fileabs(cpath string) string {
	if filepath.IsAbs(cpath) {
		return cpath
	}
	return filepath.Join(Appath, cpath)
}

// http统一验证 中间件
func HttpServerMiddleware(next http.Handler, callback func(w http.ResponseWriter, r *http.Request) bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r != nil && r.Body != nil {
				_ = r.Body.Close()
			}
		}()
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("server", "vs/1.0")
		if callback != nil && !callback(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RespJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// SetReturn 设置 返回函数Play
func SetReturn(w http.ResponseWriter, code int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	Checkerr(json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"msg":  msg,
	}), "SetReturn")
}

func SetReturnData(w http.ResponseWriter, code int, msg any, data any) {
	w.Header().Set("Content-Type", "application/json")
	Checkerr(json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}), "SetReturnData")
}

func SetReturnError(w http.ResponseWriter, r *http.Request, err error, code int, msg ...any) {
	ReqError(r, append([]any{err}, msg...)...)
	w.Header().Set("Content-Type", "application/json")
	Checkerr(json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"msg":  msg,
	}), "SetReturnError")
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

// 检查请求是否post
func CheckReqPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		Forbidden(w, "Forbidden")
		return false
	}
	return true
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

func CloseSftp(conn *sftp.Client) {
	if conn != nil {
		_ = conn.Close()
	}
}

// github.com/garyburd/redigo/redis
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

func CloseMysql(conn *sql.DB) {
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

// Getip 获取客户端IP
func Getip(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("HTTP_CLIENT_IP"); ip != "" {
		remoteAddr = ip
	} else if ip := r.Header.Get("HTTP_X_FORWARDED_FOR"); ip != "" {
		remoteAddr = ip
	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
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
