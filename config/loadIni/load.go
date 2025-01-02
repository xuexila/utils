package loadIni

import (
	"github.com/helays/utils/config"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/tools"
	"gopkg.in/ini.v1"
)

func LoadIni(i any) {
	ulogs.DieCheckerr(LoadIniBase(i), "载入配置文件失败")
}

// LoadIniBase 载入配置基础功能
func LoadIniBase(i any) error {
	return ini.MapTo(i, tools.Fileabs(config.Cpath))
}
