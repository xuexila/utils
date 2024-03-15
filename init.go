package utils

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

func init() {
	Appath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error("当前路径获取失败...", err.Error())
		os.Exit(1)
	}
}

var (
	EnableParseParamsLog = true
)

// Parseparams 解析启动参数
func Parseparams(f func()) {
	// 解析参数
	var (
		vers bool
	)
	flag.BoolVar(&Help, "h", false, "参数说明")
	flag.StringVar(&Cpath, "c", "conf.ini", "配置文件")
	flag.BoolVar(&Dbg, "debug", false, "Debug 模式")
	flag.BoolVar(&vers, "version", false, "查看版本")
	if f != nil {
		f() // 自定义
	}
	flag.Parse()
	if vers {
		fmt.Println(os.Args[0], Version, BuildTime)
		os.Exit(1)
	}
	if Help {
		flag.Usage()
		os.Exit(0)
	}
	if EnableParseParamsLog {
		Log("运行参数解析完成...")
	}

}

func LoadIni(i interface{}) {
	if err := ini.MapTo(i, Fileabs(Cpath)); err != nil {
		Error("载入配置文件错误", err.Error())
		os.Exit(1)
	}
}
