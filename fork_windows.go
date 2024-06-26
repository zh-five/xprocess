//go:build windows
// +build windows

package xprocess

import (
	"os/exec"
	"syscall"
)

func forkProcAttr(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true // 禁止出现命令行窗口
}
