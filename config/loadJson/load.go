package loadJson

import (
	"encoding/json"
	"github.com/xuexila/utils/config"
	"github.com/xuexila/utils/tools"
	"github.com/xuexila/utils/ulogs"
	"os"
)

func LoadJson(i any) {
	reader, err := os.Open(tools.Fileabs(config.Cpath))
	ulogs.DieCheckerr(err, "打开配置文件失败")
	y := json.NewDecoder(reader)
	ulogs.DieCheckerr(y.Decode(i), "解析配置文件失败")
}
