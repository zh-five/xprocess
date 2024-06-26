//go:build windows
// +build windows

package xprocess

import (
	"time"

	"golang.org/x/sys/windows"
)

func waitForEvent(event windows.Handle, onKill func()) {
	defer windows.CloseHandle(event)

	for {
		result, err := windows.WaitForSingleObject(event, windows.INFINITE)
		if err != nil {
			println("Error waiting for event:")
			panic(err)
		}
		switch result {
		case windows.WAIT_OBJECT_0:
			onKill()
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

// 返回进程是否唯一
func uniqueCheckAndKillOld(flag string, onKill func()) {
	event, name := uniOpenEvent(flag)

	// 尝试通知旧进程退出
	if event != 0 {
		err := windows.SetEvent(event) // 通知旧进程退出
		if err != nil {
			println("windows.SetEvent() error")
			panic(err)
		}
		windows.CloseHandle(event)
	}

	// 重新创建事件,用于监听
	event, err := windows.CreateEvent(nil, 0, 0, name)
	if event == 0 {
		println("windows.CreateEvent() error")
		panic(err)
	}

	go waitForEvent(event, onKill)
}

func uniOpenEvent(flag string) (event windows.Handle, name *uint16) {
	name = windows.StringToUTF16Ptr("xprocess_" + flag)
	event, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, name)
	if err != nil {
		println(err)
		event = 0
	}

	return
}

// 检查进程是不是唯一的
func uniqueCheck(flag string) (isUniq bool) {
	event, name := uniOpenEvent(flag)

	// 事件打开成功, 表示有旧进程在监听事件
	if event != 0 {
		windows.CloseHandle(event)
		isUniq = false
	} else {
		isUniq = true
	}

	// 重新创建事件,用于监听
	event, err := windows.CreateEvent(nil, 0, 0, name)
	if event == 0 {
		println("windows.CreateEvent() error")
		panic(err)
	}

	go waitForEvent(event, func() {})

	return
}
