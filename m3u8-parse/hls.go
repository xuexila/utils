package m3u8_parse

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	KeyPreg    = regexp.MustCompile(`EXT\-X\-KEY.*?URI="(.+?)"`)             // 识别 m3u8是否有key
	StreamPreg = regexp.MustCompile(`EXT\-X\-STREAM\-INF.+?BANDWIDTH=(\d+)`) // 多码率适配
)

type M3u8 struct {
	Ts      []string // ts列表
	M3u8    string   // m3u8 地址
	Key     string   // 如果有，ts内容需要用这个来解密
	Content []byte   // m3u8 文件内容
}

func M3u8Parse(m *bufio.Reader) (M3u8, error) {
	var (
		bits []int
		list []string
		res  M3u8
	)

	for {
		by, _, err := m.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}
		res.Content = append(res.Content, by...)
		res.Content = append(res.Content, []byte("\n")...)

		str := strings.TrimSpace(string(by))
		if str == "" {
			continue
		} else if strings.HasPrefix(str, "#EXT-X-STREAM-INF") {
			// 如果有多码率 需要匹配 最大的码率
			res := StreamPreg.FindStringSubmatch(str)
			if len(res) != 2 {
				continue
			}
			t, _ := strconv.Atoi(res[1])
			bits = append(bits, t)
			continue
		} else if strings.HasPrefix(str, "#EXT-X-KEY") {
			// 如果匹配到 key 就获取key
			keyresult := KeyPreg.FindStringSubmatch(str)
			if len(keyresult) != 2 {
				return res, errors.New("解析ts 密钥失败")
			}
			res.Key = string(keyresult[1])
			continue
		} else if strings.HasPrefix(str, "#") {
			continue
		}
		list = append(list, str)
	}

	if len(bits) < 1 || bits == nil {
		if list == nil {
			return res, errors.New("m3u8地址识别失败")
		}
		res.Ts = list
		return res, nil
	}
	if len(bits) != len(list) || list == nil {
		return res, errors.New("码流识别失败")
	}
	var swap []int
	copy(swap, bits)
	sort.Ints(swap)
	i := bits[len(bits)-1]
	for index, value := range bits {
		if value == i {
			res.M3u8 = list[index]
			break
		}
	}
	return res, nil
}
