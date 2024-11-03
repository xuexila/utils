package config

import (
	"regexp"
	"time"
)

var (
	Help   bool   // 打印显示帮助信息
	Cpath  string // 配置文件路径
	Appath string // 当前路径
	Dbg    bool   // Debug 模式

	CstSh = time.FixedZone("CST", 8*3600) // 东八区

	PublicKeyByt  []byte // 公钥
	PrivateKeyByt []byte // 私钥

	SpecialChartPreg = regexp.MustCompile(`[\s;!@#$%^&*()\[\]\:\"\']`)
	WrapSplit        = regexp.MustCompile("\r?\n")
	PhoneRegex       = regexp.MustCompile(`^1(3|4|5|6|7|8|9)\d{9}`)                       // 手机号验证正则
	EmailRegex       = regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,6}$`) // 邮箱验证正则

	PageCharacterSetPreg = regexp.MustCompile(`(?i)<meta[\s\S]+?charset=['"]?([\w-]+)['"][\s\S]*?>`)              // 匹配页面 字符集
	ALinkRegex           = regexp.MustCompile(`(?i)\<(?:a|mip\-link)\s[\s\S]*?href=['"]([\s\S]*?)['"][\s\S]*?\>`) // 匹配页面 所有a标签的正则表达式
	Nofollow             = regexp.MustCompile(`(?i)rel=['"]nofollow['"]`)                                         // 匹配 rel=nofollow
	TitleRegex           = regexp.MustCompile(`(?i)\<title[\s\S]*?\>([\s\S]*?)\<\/title\>`)                       // 匹配页面 title的正则表达式

	EnableHttpserver     bool
	CloseHttpserverSig   = make(chan byte, 1)
	EnableParseParamsLog = true
)
