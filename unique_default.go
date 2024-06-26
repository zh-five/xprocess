//go:build !windows && !plan9
// +build !windows,!plan9

package xprocess

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

func uniIsProcessRunning(pid int) bool {
	// syscall.Kill函数在这里被用来向进程发送信号0
	// 成功发送信号表示进程存在，返回nil错误
	// 失败（即错误不为nil）则可能表示进程不存在

	return syscall.Kill(pid, 0) == nil
}

var uniPidFile = ""

// 检查进程是不是唯一的，不是则kill旧进程
// 为避免误杀和优雅退出, 此处kill不是强杀. 需在 onKill() 方法中自行处理退出操作
func uniqueCheckAndKillOld(flag string, onKill func()) {
	uniqueCheckAndDo(flag, func(proc *os.Process) {
		err := proc.Signal(syscall.SIGHUP) // 通知旧进程退出
		fmt.Println("proc.Signal(syscall.SIGHUP):", err)
	})

	go func() {
		<-NewSignal(syscall.SIGHUP).Done()
		onKill()
	}()
}

// 检查进程是不是唯一的
func uniqueCheck(flag string) bool {
	uniSetPidFile(flag)

	pid := uniReadPid()
	if pid > 0 {
		if uniIsProcessRunning(pid) {
			return false
		}
	}

	uniSavePid()
	return true
}

func uniqueCheckAndDo(flag string, do func(proc *os.Process)) {
	uniSetPidFile(flag)

	pid := uniReadPid()
	if pid > 0 {
		if uniIsProcessRunning(pid) {
			proc, err := os.FindProcess(pid)
			if err == nil {
				do(proc)
			}
		}
	}

	uniSavePid()
}

func uniReadPid() int {
	str, err := os.ReadFile(uniPidFile)
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		panic(err)
	}

	pid, err := strconv.Atoi(string(str))
	if err != nil {
		panic(err) // pid 文件格式异常,可能被篡改了
	}
	return pid
}

func uniSavePid() {
	str := strconv.Itoa(os.Getpid())
	err := os.WriteFile(uniPidFile, []byte(str), 0o600)
	if err != nil {
		panic(err)
	}
}

// 获取
var uniGetPidFile = func(flag string) (pidFile string) {
	dir, err := os.UserCacheDir()
	defer func() {
		name := fmt.Sprintf("xprocess_%s.pid", flag)
		pidFile = filepath.Join(dir, name)
	}()

	if err == nil && uniCheckDirRW(dir) {
		return
	}

	dir, err = os.UserConfigDir()
	if err == nil && uniCheckDirRW(dir) {
		return
	}

	dir, err = os.UserHomeDir()
	if err == nil && uniCheckDirRW(dir) {
		return
	}

	dir = os.TempDir()
	if !uniCheckDirRW(dir) {
		panic("无法找到可以存储pid的目录")
	}
	return
}

func uniSetPidFile(flag string) {
	if uniPidFile != "" {
		return
	}

	uniPidFile = uniGetPidFile(flag)
	fmt.Println("pid file:", uniPidFile)
}

// 检查目录是否有读写权限
func uniCheckDirRW(dirPath string) bool {
	tempFile, err := os.CreateTemp(dirPath, "test_")
	if err != nil {
		return false
	}

	tempFile.Close()
	err = os.Remove(tempFile.Name())

	return err == nil
}
