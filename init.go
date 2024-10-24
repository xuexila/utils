package utils

import (
	"github.com/xuexila/utils/config"
	"github.com/xuexila/utils/ulogs"
	"os"
	"path/filepath"
)

func init() {
	var err error
	config.Appath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		ulogs.Error("当前路径获取失败...", err.Error())
		os.Exit(1)
	}
}
