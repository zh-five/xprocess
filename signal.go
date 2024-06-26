package xprocess

import (
	"os"
	"os/signal"
	"syscall"
)

type Signal struct {
	sigs []os.Signal
}

// 创建一个监听
func NewSignal(sigs ...os.Signal) *Signal {
	return &Signal{sigs: sigs}
}

// 无条件结束程序(不能被捕获、阻塞或忽略)
func (sf *Signal) SignalKill() *Signal {
	sf.sigs = append(sf.sigs, syscall.SIGKILL)
	return sf
}

// 用户发送INTR字符(Ctrl+C)触发
func (sf *Signal) SignalCtrlC() *Signal {
	sf.sigs = append(sf.sigs, syscall.SIGINT)
	return sf
}

// 结束程序(可以被捕获、阻塞或忽略)
func (sf *Signal) SignalEnd() *Signal {
	sf.sigs = append(sf.sigs, syscall.SIGTERM)
	return sf
}

// 程序退出
func (sf *Signal) SignalAllExit() *Signal {
	sf.sigs = append(sf.sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	return sf
}

// 返回的通道可一直等待和读取出现的信号(限定为所监听的信号)
func (sf *Signal) Notify() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sf.sigs...)

	return ch
}

// 出现任何一个监听的信号, 将关闭返回的通道
func (sf *Signal) Done() chan struct{} {
	ch := make(chan struct{})
	go func() {
		chN := sf.Notify()
		<-chN
		signal.Stop(chN)
		close(ch)
	}()

	return ch
}
