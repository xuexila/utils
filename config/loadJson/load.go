package loadJson

import (
	"encoding/json"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/config"
	"github.com/helays/utils/tools"
	"github.com/helays/utils/ulogs"
	"os"
)

func LoadJson(i any) {
	reader, err := os.Open(tools.Fileabs(config.Cpath))
	defer osClose.CloseFile(reader)
	ulogs.DieCheckerr(err, "打开配置文件失败")
	y := json.NewDecoder(reader)
	ulogs.DieCheckerr(y.Decode(i), "解析配置文件失败")
}
