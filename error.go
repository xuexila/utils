package utils

import (
	"log"
	"os"
)

// Recover 捕获系统异常
func Recover() {
	if r := recover(); r != nil {
		Error("系统异常，捕获结果", r)
	}
}

// Log 打印正确日志。
func Log(i ...interface{}) {
	// log.SetPrefix("[用户日志]")
	log.SetOutput(os.Stdout)
	log.Println(i...)
}

func Debug(i ...interface{}) {
	if Dbg {
		Log("[debug]", i)
	}
}

// Checkerr 检查错误
func Checkerr(err error, i ...interface{}) {
	if err == nil {
		return
	}
	Error(i, err)
}

// DieCheckerr 检查错误，打印并输出错误信息
func DieCheckerr(err error, i ...interface{}) {
	if err == nil {
		return
	}
	Error(i, err)
	os.Exit(1)
}

// ReturnCheckerr 检查错误，有异常就返回false
func ReturnCheckerr(err error, i ...interface{}) bool {
	if err == nil {
		return true
	}
	Error(i, err)
	return false
}

func ErrorReturn(i ...interface{}) bool {
	Error(i...)
	return false
}

// Error 打印错误信息
func Error(i ...interface{}) {
	log.SetPrefix("")
	log.SetOutput(os.Stderr)
	//_lst:=i[len(i)-1]
	//fmt.Println("ss",_lst==nil)
	log.Println(i...)
}

func Pfunc(a ...interface{}) {
	// log.SetPrefix("[用户异常]")
	log.SetOutput(os.Stdout)
	log.Println(a...)
}
