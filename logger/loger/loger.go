package loger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

/**
日志系统初始化流程
1、判断 各种日志是否有限制单个文件的最大数量，单位 mb
2、判断 各种日志是否需要保留指定数量的文件数
3、判断 是否需要将错误日志输出到正常日志文件中
4、如果日志文件不限制文件大小，保留份数不生效，如果将错误日志输出到正常日志文件中，错误日志文件数量，保留份数不生效
*/

var (
	lock sync.Mutex
)

func Init(l Loger) (*Loger, error) {
	if l.StdoutFile == "" {
		return nil, errors.New("正常日志输出目录空")
	}

	if !l.RedirectStderr && l.StderrFile == "" {
		return nil, errors.New("错误日志输出目录空")
	}

	if l.StdoutMaxBytes >= 1 {
		l.StdoutMaxBytes = l.StdoutMaxBytes * 1024 * 1024
	}
	if l.StderrMaxBytes >= 1 {
		l.StderrMaxBytes = l.StderrMaxBytes * 1024 * 1024
	}
	if l.StdoutBackups < 1 {
		l.StdoutBackups = 1
	}
	if l.StderrBackups < 1 {
		l.StderrBackups = 1
	}
	if l.RedirectStderr {
		l.StderrFile = l.StdoutFile
		l.StderrBackups = l.StdoutBackups
	}
	// 创建日志文件保存目录
	if err := mkdir(filepath.Dir(l.StdoutFile)); err != nil {
		return nil, err
	}
	if err := mkdir(filepath.Dir(l.StderrFile)); err != nil {
		return nil, err
	}
	var lo = new(Loger)
	lo.Name = l.Name
	lo.RedirectStderr = l.RedirectStderr
	lo.StdoutFile = l.StdoutFile
	lo.StdoutMaxBytes = l.StdoutMaxBytes
	lo.StdoutBackups = l.StdoutBackups
	lo.StderrFile = l.StderrFile
	lo.StderrMaxBytes = l.StderrMaxBytes
	lo.StderrBackups = l.StderrBackups
	lo.wsig = make(chan byte)
	// 打开文件
	if err := lo.openStdoutFile(); err != nil {
		return nil, err
	}
	if err := lo.openStderrFile(); err != nil {
		return nil, err
	}
	if lo.Name == "" {
		lo.Name = "Loger"
	}
	go lo.heartMiddleware()
	return lo, nil
}

// Log 记录正常日志
func (l *Loger) Log(i ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if err := l.openStdoutFile(); err != nil {
		fmt.Println("Loger", "Log", err)
		return
	}
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[" + l.Name + "][Log] ")
	log.SetOutput(l.stdoutFile)
	log.Println(i...)
	_ = l.stdoutFile.Sync()
	l.wsig <- 1

}

// 记录错误日志
func (l *Loger) Error(i ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if err := l.openStderrFile(); err != nil {
		fmt.Println("Loger", "Error", err)
		return
	}
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[" + l.Name + "][Error] ")
	log.SetOutput(l.stderrFile)
	log.Println(i...)
	_ = l.stderrFile.Sync()
	l.wsig <- 2
}

// 心跳部分，用于检测 文件是否需要处理，比如分割，重命名等。
func (l *Loger) heartMiddleware() {
	for {

		switch <-l.wsig {
		case 1:
			if l.StdoutMaxBytes == 0 {
				continue
			}
			l.splitFile(1)
			// 正常日志 需要定时分割
		case 2:
			if (l.RedirectStderr && l.StdoutMaxBytes == 0) || (!l.RedirectStderr && l.StderrMaxBytes == 0) {
				// 错误日志 重定向到 正常日志，正常日志不限制文件大小，跳过处理
				// 错误日志不重定向，错误日志不限制文件大小，跳过处理
				continue
			} else if l.RedirectStderr {
				// 判断分割 正常日志
				l.splitFile(1)
				continue
			}
			// 判断分割错误日志
			l.splitFile(2)
		}
	}
}

func (l *Loger) splitFile(t int) {
	var fsize int64
	if t == 1 || l.RedirectStderr {
		if err := l.openStdoutFile(); err != nil {
			fmt.Println("Loger", "splitFile", `t == 1 || l.RedirectStderr`, "l.openStdoutFile()", l.StdoutFile, err)
		} else if stat, err := l.stdoutFile.Stat(); err != nil {
			fmt.Println("Loger", "splitFile", `t == 1 || l.RedirectStderr`, "stdoutFile.Stat()", l.StdoutFile, err)
		} else {
			fsize = stat.Size()
		}
		if fsize < l.StdoutMaxBytes {
			return
		}
	} else {
		if err := l.openStderrFile(); err != nil {
			fmt.Println("Loger", "splitFile", "l.openStderrFile()", l.StderrFile, err)
		} else if stat, err := l.stderrFile.Stat(); err != nil {
			fmt.Println("Loger", "splitFile", "stderrFile.Stat()", l.StderrFile, err)
		} else {
			fsize = stat.Size()
		}
		if fsize < l.StderrMaxBytes {
			return
		}
	}

	// 这里是需要分割的。
	if t == 1 && l.StdoutBackups == 1 {
		// 判断正常日志
		if err := l.openStdoutFile(); err != nil {
			fmt.Println("Loger", "splitFile", "l.openStdoutFile()", l.StdoutFile, err)
		} else if err := l.stdoutFile.Truncate(0); err != nil {
			fmt.Println("Loger", "splitFile", "stdoutFile.Truncate(0)", l.StdoutFile, err)
		}
	} else if l.RedirectStderr && l.StdoutBackups == 1 {
		// 错误日志需要重定向，并且 正常日志只保留一份
		if err := l.openStderrFile(); err != nil {
			fmt.Println("Loger", "splitFile", "l.openStderrFile()", l.StdoutFile, err)
		} else if err := l.stderrFile.Truncate(0); err != nil {
			fmt.Println("Loger", "splitFile", "stderrFile.Truncate(0)", l.StdoutFile, err)
		} else if err := l.openStdoutFile(); err != nil {
			fmt.Println("Loger", "splitFile", "l.openStdoutFile()", l.StdoutFile, err)
		} else if _, err := l.stdoutFile.Seek(0, 0); err != nil {
			fmt.Println("Loger", "splitFile", "stdoutFile.Seek(0,0)", l.StdoutFile, err)
		}
	} else if t == 2 && l.StderrBackups == 1 {
		// 错误日志 清空
		if err := l.openStderrFile(); err != nil {
			fmt.Println("Loger", "splitFile", "l.openStderrFile()", l.StderrFile, err)
		} else if err := l.stderrFile.Truncate(0); err != nil {
			fmt.Println("Loger", "splitFile", "stderrFile.Truncate(0)", l.StderrFile, err)
		}
	} else if t == 1 || l.RedirectStderr {
		// 正常日志，需要保留多份文件的
		// 错误日志需要重定向到 正常日志
		dir := filepath.Dir(l.StdoutFile)
		filename := filepath.Base(l.StdoutFile)
		l.updateFileName(dir, filename, l.StdoutBackups)

		_ = l.stdoutFile.Close()
		l.stdoutFile = nil
		if l.RedirectStderr {
			_ = l.stderrFile.Close()
			l.stderrFile = nil
		}
	} else if t == 2 {
		// 错误日志 需要保留多份日志文件
		dir := filepath.Dir(l.StderrFile)
		filename := filepath.Base(l.StderrFile)
		l.updateFileName(dir, filename, l.StderrBackups)
		_ = l.stderrFile.Close()
		l.stderrFile = nil
	}
}

// 更新日志文件名次
func (l *Loger) updateFileName(dir, filename string, backups int) {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Loger", "updateFileName", "ioutil.ReadDir(dir)", dir, err)
		return
	}

	numArr := []int{}

	for _, item := range dirs {
		if item.IsDir() {
			continue
		}
		fname := item.Name()
		if fname == filename {
			continue
		}
		if !strings.HasPrefix(fname, filename) {
			continue
		}
		stmp := strings.Split(fname, ".")
		sufId := stmp[len(stmp)-1]
		sId, err := strconv.Atoi(sufId)
		if err != nil {
			continue
		}
		oldName := filepath.Join(dir, fname)
		if sId >= (backups - 1) {
			_ = os.RemoveAll(oldName)
			continue
		}
		numArr = append(numArr, sId)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(numArr))) //对文件后缀进行倒序排列

	for _, id := range numArr {
		oldName := filepath.Join(dir, filename+"."+strconv.Itoa(id))
		newName := filepath.Join(dir, filename+"."+strconv.Itoa(id+1))
		_ = os.Rename(oldName, newName)
	}
	oldName := filepath.Join(dir, filename)
	_ = os.Rename(oldName, oldName+".1")
}

// 打开正常日志文件句柄
func (l *Loger) openStdoutFile() (err error) {
	if l.stdoutFile != nil {
		return nil
	}

	l.stdoutFile, err = os.OpenFile(l.StdoutFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	return
}

// 打开错误日志文件句柄
func (l *Loger) openStderrFile() (err error) {
	if l.stderrFile != nil {
		return nil
	}
	l.stderrFile, err = os.OpenFile(l.StderrFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	return
}
