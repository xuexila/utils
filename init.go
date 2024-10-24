package utils

import (
	"flag"
	"fmt"
	"github.com/xuexila/utils/config"
	"github.com/xuexila/utils/ulogs"
	"os"
	"path/filepath"
)

func init() {
	var err error
	config.Appath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		ulogs.Error("当前路径获取失败...", err.Error())
		os.Exit(1)
	}
}

// Parseparams 解析启动参数
func Parseparams(f func()) {
	// 解析参数
	var (
		vers bool
	)
	flag.BoolVar(&config.Help, "h", false, "参数说明")
	flag.StringVar(&config.Cpath, "c", "conf.ini", "配置文件")
	flag.BoolVar(&config.Dbg, "debug", false, "Debug 模式")
	flag.BoolVar(&vers, "version", false, "查看版本")
	if f != nil {
		f() // 自定义
	}
	flag.Parse()
	if vers {
		fmt.Println(os.Args[0], Version, BuildTime)
		os.Exit(1)
	}
	if config.Help {
		flag.Usage()
		os.Exit(0)
	}
	if config.EnableParseParamsLog {
		ulogs.Log("运行参数解析完成...")
	}
}
