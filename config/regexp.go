package config

import "regexp"

var (
	NumberRegexp = regexp.MustCompile(`^[0-9]+$`)

	SpecialChartPreg = regexp.MustCompile(`[\s;!@#$%^&*()\[\]\:\"\']`)
	WrapSplit        = regexp.MustCompile("\r?\n")
	PhoneRegex       = regexp.MustCompile(`^1(3|4|5|6|7|8|9)\d{9}`)                       // 手机号验证正则
	EmailRegex       = regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,6}$`) // 邮箱验证正则

	PageCharacterSetPreg = regexp.MustCompile(`(?i)<meta[\s\S]+?charset=['"]?([\w-]+)['"][\s\S]*?>`)              // 匹配页面 字符集
	ALinkRegex           = regexp.MustCompile(`(?i)\<(?:a|mip\-link)\s[\s\S]*?href=['"]([\s\S]*?)['"][\s\S]*?\>`) // 匹配页面 所有a标签的正则表达式
	Nofollow             = regexp.MustCompile(`(?i)rel=['"]nofollow['"]`)                                         // 匹配 rel=nofollow
	TitleRegex           = regexp.MustCompile(`(?i)\<title[\s\S]*?\>([\s\S]*?)\<\/title\>`)                       // 匹配页面 title的正则表达式
)
