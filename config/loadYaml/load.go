package loadYaml

import (
	"fmt"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/config"
	"github.com/helays/utils/tools"
	"github.com/helays/utils/ulogs"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadYaml(i any) {
	ulogs.DieCheckerr(LoadYamlBase(i), "解析配置文件失败")
}

func LoadYamlBase(i any) error {
	reader, err := os.Open(tools.Fileabs(config.Cpath))
	defer osClose.CloseFile(reader)
	if err != nil {
		return fmt.Errorf("打开配置文件失败：%s", err.Error())
	}
	y := yaml.NewDecoder(reader)
	return y.Decode(i)
}
