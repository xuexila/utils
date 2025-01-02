package loadJson

import (
	"encoding/json"
	"fmt"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/config"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/tools"
	"os"
)

func LoadJson(i any) {
	ulogs.DieCheckerr(LoadJsonBase(i), "解析配置文件失败")
}

func LoadJsonBase(i any) error {
	reader, err := os.Open(tools.Fileabs(config.Cpath))
	defer osClose.CloseFile(reader)
	if err != nil {
		return fmt.Errorf("打开配置文件失败：%s", err.Error())
	}
	y := json.NewDecoder(reader)
	return y.Decode(i)
}
