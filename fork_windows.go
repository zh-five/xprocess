//go:build windows
// +build windows

package xprocess

import (
	"os/exec"
)

func forkProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr.HideWindow = true // 禁止出现命令行窗口
}
