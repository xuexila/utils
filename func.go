package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	url2 "net/url"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// SnakeString 驼峰转蛇形
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// CamelString 蛇形转驼峰
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

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
	if EnableHttpserver {
		CloseHttpserverSig <- 1
		_ = <-CloseHttpserverSig
	}
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
		undatas, err := io.ReadAll(r)
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
	if o == nil {
		return false
	}
	s = strings.TrimSpace(s)
	for _, i := range o {
		i = strings.TrimSpace(i)
		if i == s {
			return true
		}
	}
	return false
}

// SearchIntSlice 在整数切片中搜索指定的元素，并返回是否找到。
// 参数:
//
//	s - 待搜索的整数。
//	arr - 整数切片，将被搜索。
//
// 返回值:
//
//	如果找到 s 在 arr 中，则返回 true；否则返回 false。
func SearchIntSlice(s int, arr []int) bool {
	if arr == nil {
		return false
	}
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

func SearchInt64Slice(s int64, arr []int64) bool {
	if arr == nil {
		return false
	}
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

// MinMaxAvgSum 获取数组中最大值最小值平均值和求和
func MinMaxAvgSum(nums []int) (min int, max int, avg float64, sum int) {
	if len(nums) == 0 {
		return 0, 0, 0, 0
	}
	min, max, sum = nums[0], nums[0], nums[0]
	for _, num := range nums[1:] {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
		sum += num
	}
	avg = float64(sum) / float64(len(nums))
	return
}

// AnySlice2Str 将任意切片转成字符串
func AnySlice2Str(slice []any, _sep ...string) string {
	var builder strings.Builder
	l := len(slice)
	sep := ","
	if len(_sep) > 0 {
		sep = _sep[0]
	}
	for index, elem := range slice {
		// 使用 fmt.Sprint 将任何类型转换为字符串形式
		strElem := fmt.Sprint(elem)
		if strElem == "" {
			continue
		}
		builder.WriteString(strElem)
		// 可以选择在此处添加分隔符，如空格、逗号等
		if index < (l - 1) {
			builder.WriteString(sep)
		}
	}

	return builder.String()
}

type test struct {
	uu useTest
	useTest
}

type useTest struct {
	a string
}

// Map2Struct 将map转换为结构体
// dst 需要传入一个变量的指针
func Map2Struct(dst any, src map[string]any, customConvert map[string]func(dst any, src map[string]any) error) error {
	// 这里通过反射，将map转换为结构体
	val := reflect.ValueOf(dst).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		// 如果这个字段是一个匿名结构体，还需要递归处理
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			// 递归处理嵌套结构体，注意这里不能使用field
			if err := Map2Struct(val.Field(i).Addr().Interface(), src, customConvert); err != nil {
				return err
			}
			continue
		}
		// 如果有自定义转换函数，则使用自定义转换函数
		if f, ok := customConvert[field.Name]; ok {
			if err = f(dst, src); err != nil {
				return fmt.Errorf("自定义转换函数%s执行失败：%v", typ.Field(i).Name, err)
			}
			continue
		}
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		jsonTag=strings.Split(jsonTag,",")[0]
		// 获取结构体json标签
		// 检查 map 中是否存在对应的键
		value, ok := src[jsonTag]
		if !ok {
			continue
		}
		// 设置字段的值
		fieldVal := val.Field(i)
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(fmt.Sprintf("%v", value))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			_tv, err := Any2int(value)
			if err != nil {
				return fmt.Errorf("字段%s转int失败：%v", jsonTag, err)
			}
			fieldVal.SetInt(_tv)
		case reflect.Float32, reflect.Float64:
			tv, err := Any2float64(value)
			if err != nil {
				return fmt.Errorf("字段 %s 转 float64 失败: %v", jsonTag, err)
			}
			fieldVal.SetFloat(tv)
		case reflect.Bool:
			tv, err := Any2bool(value)
			if err != nil {
				return fmt.Errorf("字段 %s 转 bool 失败: %v", jsonTag, err)
			}
			fieldVal.SetBool(tv)
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Uint8 { // []byte
				tv, err := Any2bytes(value)
				if err != nil {
					return fmt.Errorf("字段 %s 转 []byte 失败: %v", jsonTag, err)
				}
				fieldVal.SetBytes(tv)
			}
		default:
		}
	}
	return nil
}

func Any2string(v any) string {
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case int32:
		return strconv.FormatInt(int64(v.(int32)), 10)
	case float32:
		return Float32tostring(v.(float32))
	case float64:
		return Float64tostring(v.(float64))
	}
	return fmt.Sprintf("%v", v)
}

// Any2int 尝试将任意类型转换为 int
func Any2int(v any) (int64, error) {
	switch v := v.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		if v == "" {
			return 0, nil
		}
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("无法将类型 %T 转换为 int", v)
	}
}

// Any2float64 尝试将任意类型转换为 float64
func Any2float64(v any) (float64, error) {
	switch v := v.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0.0, fmt.Errorf("无法将类型 %T 转换为 float64", v)
	}
}

// Any2bool 尝试将任意类型转换为 bool
func Any2bool(v any) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("无法将类型 %T 转换为 bool", v)
	}
}

// Any2bytes 尝试将任意类型转换为 []byte
func Any2bytes(v any) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("无法将类型 %T 转换为 []byte", v)
	}
}

// GetIpVersion 解析ip地址，确认ip版本
func GetIpVersion(ip string) (string, error) {
	_ip := net.ParseIP(ip)
	if _ip == nil {
		return "", errors.New("ip地址不合法")
	}
	if _ip.To4() != nil {
		return "ipv4", nil
	}
	return "ipv6", nil
}

// Str2StrSlice 字符串转切片
func Str2StrSlice(values string) ([]string, error) {
	values = strings.TrimSpace(values)
	if values == "" {
		return nil, nil
	}
	var slice []string
	if err := json.Unmarshal([]byte(values), &slice); err != nil {
		return nil, fmt.Errorf("解析数据 [%s] 失败：%v", values, err)
	}
	return slice, nil
}
