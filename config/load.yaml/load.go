package load_yaml

import (
	"gitlab.itestor.com/helei/utils.git"
	"gitlab.itestor.com/helei/utils.git/config"
	"gitlab.itestor.com/helei/utils.git/ulogs"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadYaml(i any) {
	reader, err := os.Open(utils.Fileabs(config.Cpath))
	ulogs.DieCheckerr(err, "打开配置文件失败")
	y := yaml.NewDecoder(reader)
	ulogs.DieCheckerr(y.Decode(i), "解析配置文件失败")
}
