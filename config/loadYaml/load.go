package loadYaml

import (
	"github.com/helays/utils/config"
	"github.com/helays/utils/tools"
	"github.com/helays/utils/ulogs"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadYaml(i any) {
	reader, err := os.Open(tools.Fileabs(config.Cpath))
	ulogs.DieCheckerr(err, "打开配置文件失败")
	y := yaml.NewDecoder(reader)
	ulogs.DieCheckerr(y.Decode(i), "解析配置文件失败")
}
