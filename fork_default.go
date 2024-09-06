//go:build !windows && !plan9
// +build !windows,!plan9

package xprocess

import (
	"os/exec"
	"syscall"
)

func CmdProcAttr(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setsid = true // 可脱离父进程独立运行，不受父进程退出的影响
}
