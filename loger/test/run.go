package main

import (
	"fmt"
	"github.com/helays/utils/config/loadIni"
	"github.com/helays/utils/config/parseCmd"
	"github.com/helays/utils/loger"
	"strings"
	"time"
)

type config struct {
	Log loger.Loger
}

func main() {
	var log = new(config)
	parseCmd.Parseparams(nil)
	loadIni.LoadIni(log)

	loger.Init(log.Log)
	if err := log.Log.Init(); err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			log.Log.Error(time.Now().Unix())
		}
	}()
	for {
		log.Log.Log(strings.Repeat(time.Now().String(), 2))
	}
}
