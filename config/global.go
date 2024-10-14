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

	SpecialChartPreg     = regexp.MustCompile(`[\s;!@#$%^&*()\[\]\:\"\']`)
	EnableHttpserver     bool
	CloseHttpserverSig   = make(chan byte, 1)
	EnableParseParamsLog = true
)
