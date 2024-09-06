//go:build windows
// +build windows

package xprocess

import (
	"os/exec"
	"syscall"
)

func CmdProcAttr(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true // 禁止出现命令行窗口

	// 脱离主进程运行
	const DETACHED_PROCESS = 0x00000008
	cmd.SysProcAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP | DETACHED_PROCESS
}
