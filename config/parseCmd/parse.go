package parseCmd

import (
	"flag"
	"fmt"
	"github.com/helays/utils"
	"github.com/helays/utils/config"
	"github.com/helays/utils/ulogs"
	"os"
)

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
		fmt.Println(os.Args[0], utils.Version, utils.BuildTime)
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
