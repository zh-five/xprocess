package xprocess

import (
	"context"
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

// 监视事件, 事件由通道返回
func (sf *Signal) Notify() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sf.sigs...)

	return ch
}

// 发生任何监视的事件时, cancel()
func (sf *Signal) WithCancel(ctx context.Context) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sf.sigs...)
	go func() {
		select {
		case <-ctx.Done():
			signal.Stop(ch)
		case <-ch:
			cancel()
			signal.Stop(ch)
		}
	}()

	return ctx, cancel
}
