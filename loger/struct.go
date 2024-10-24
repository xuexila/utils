package loger

import "os"

type Loger struct {
	Name           string `ini:"name"`             // 当前日志模块名字
	RedirectStderr bool   `ini:"redirect_stderr"`  // 将错误输出，重定向到正常日志
	StdoutFile     string `ini:"stdout_file"`      // 正常日志文件
	StdoutMaxBytes int64  `ini:"stdout_max_bytes"` // 正常日志最大字节数 单位MB
	StdoutBackups  int    `ini:"stdout_backups"`   // 正常日志保存数量
	StderrFile     string `ini:"stderr_file"`      // 错误日志文件
	StderrMaxBytes int64  `ini:"stderr_max_bytes"` // 错误日志最大字节数 单位MB
	StderrBackups  int    `ini:"stderr_backups"`   // 错误日志保存数量

	stdoutFile *os.File  // 标准正常日志
	stderrFile *os.File  // 标准错误日志
	wsig       chan byte // 1、正常日志 2、错误日志

}
