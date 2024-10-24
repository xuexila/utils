package load_ini

import (
	"github.com/xuexila/utils"
	"github.com/xuexila/utils/config"
	"github.com/xuexila/utils/ulogs"
	"gopkg.in/ini.v1"
	"os"
)

func LoadIni(i any) {
	if err := ini.MapTo(i, utils.Fileabs(config.Cpath)); err != nil {
		ulogs.Error("载入配置文件错误", err.Error())
		os.Exit(1)
	}
}
