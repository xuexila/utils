package shell

import (
	"bytes"
	"errors"
	"gitlab.itestor.com/helei/utils.git"
	"os/exec"
)

// ExecShell 执行 shell语句
func ExecShell(name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = bytes.NewBuffer(nil)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	utils.Log("shell 执行命令", cmd.String())

	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}

func ExecShellQuit(name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = bytes.NewBuffer(nil)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}

// ExecCtlShell 可手动结束的 shell命令
func ExecCtlShell(stop chan byte, name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = new(bytes.Buffer)
		end = make(chan byte)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	utils.Log("shell 执行命令", cmd.String())
	go func(cmd *exec.Cmd, stop, end chan byte) {
		for {
			select {
			case <-stop:
				if cmd == nil {
					utils.Error("分析EXEC 不存在", cmd.String())
					continue
				}
				if cmd.Process == nil {
					utils.Error("分析 Process不存在", cmd.String())
					continue
				}
				if err := cmd.Process.Kill(); err != nil {
					utils.Error("手动结束进程失败", err)
				} else {
					utils.Log("手动结束命令", cmd.String())
				}

			case <-end:
				utils.Log("shell 执行完成", cmd.String())
				return
			}
		}
	}(cmd, stop, end)
	defer func() {
		go func() {
			end <- 1
		}()

	}()
	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}
