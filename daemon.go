package xprocess

import (
	"context"
)

// 后台守护进程模式
// 1.主进程(fork后退出) fork=> 守护进程 fork=> 业务进程
// 2.守护进程永远只会启动一个, 新的启动了旧的就会退出. 唯一性依赖 flagStr
func Daemon(
	flagStr string, // 唯一标识, 用于进程间通信
	logFile string, // 守护进程和业务进程的输出会写入此文件
	isKill bool, // 是否杀掉守护进程和业务进程
	isDaemon bool, // 是否启动守护进程

) (ctx context.Context, cancel func(), isExit bool) {

	ctx, cancel = context.WithCancel(context.Background())
	isExit = true

	if isKill {
		// 杀掉管理进程, 管理进程会杀掉业务进程
		UniqueCheckAndKillOld(flagStr, func() {})
		return
	}

	if isDaemon {
		// 启动一个子进程, 作为管理进程, 用于维护和重启业务进程
		cmd, err := Fork2Log(logFile)
		if err != nil {
			panic(err)
		}
		if cmd != nil { // 主进程, 退出
			return
		}

		// 跳过业务进程
		if !IsForkPassing() {
			// 管理进程保持唯一性
			UniqueCheckAndKillOld(flagStr, func() {
				cancel()
			})
		}

		// 启动一个子进程, 并维持起常驻内存, 若异常退出则再次启动一个
		if AlwaysFork2Std(ctx) { // 主进程将阻塞在此
			// 主进程退出
			return
		}

		//子进程可以继续往后执行
	}

	isExit = false
	return
}
