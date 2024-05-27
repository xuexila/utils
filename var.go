package utils

import (
	"regexp"
	"time"
)

// 全局变量包

const (
	version string = "1.0"   // 版本号
	Salt    string = "helei" // 加密salt
)

var (
	Help      bool                            // 打印显示帮助信息
	Cpath     string                          // 配置文件路径
	Appath    string                          // 当前路径
	Dbg       bool                            // Debug 模式
	Version   string                          // 打印版本\
	BuildTime string                          // 程序编译时间
	CstSh     = time.FixedZone("CST", 8*3600) // 东八区
	err       error                           // 错误

	PublicKeyByt  []byte // 公钥
	PrivateKeyByt []byte // 私钥

	defaultLetters     = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	SpecialChartPreg   = regexp.MustCompile(`[\s;!@#$%^&*()\[\]\:\"\']`)
	EnableHttpserver   bool
	CloseHttpserverSig = make(chan byte, 1)
)
