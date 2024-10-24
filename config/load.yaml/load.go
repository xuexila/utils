package load_yaml

import (
	"github.com/xuexila/utils"
	"github.com/xuexila/utils/config"
	"github.com/xuexila/utils/ulogs"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadYaml(i any) {
	reader, err := os.Open(utils.Fileabs(config.Cpath))
	ulogs.DieCheckerr(err, "打开配置文件失败")
	y := yaml.NewDecoder(reader)
	ulogs.DieCheckerr(y.Decode(i), "解析配置文件失败")
}
