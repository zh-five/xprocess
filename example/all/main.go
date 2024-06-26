package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/zh-five/xprocess"
)

var Version = "test"

func main() {
	isKill := flag.Bool("k", false, "kill 现在的进程, 有此参数时, 忽略 -d 参数")
	isDaemon := flag.Bool("d", false, "后台运行")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())

	flagStr := "gv_cdp"
	logFile := filepath.Join(os.TempDir(), fmt.Sprintf("gv_cdp_%s.log", time.Now().Format("20060102_150405")))
	fmt.Println("log file:", logFile)

	pid := os.Getpid()
	fmt.Println("pid:", pid)

	if *isKill {
		// 杀掉管理进程, 管理进程会杀掉业务进程
		xprocess.UniqueCheckAndKillOld(flagStr, func() {})
		fmt.Println("已经向旧进程发送退出信号")
		return
	}

	if *isDaemon {
		// 启动一个子进程, 作为管理进程, 用于维护和重启业务进程
		cmd, err := xprocess.Fork2Log(logFile)
		if err != nil {
			panic(err)
		}
		if cmd != nil { // 主进程, 退出
			return
		}

		fmt.Println(pid, ">= 管理进程")
		// 跳过业务进程
		if !xprocess.IsForkPassing() {
			fmt.Println(pid, " = 管理进程: UniqueCheckAndKillOld()")
			// 管理进程保持唯一性
			xprocess.UniqueCheckAndKillOld(flagStr, func() {
				fmt.Println(pid, " = 管理进程: cancel")
				cancel()
			})
		}

		// 启动一个子进程, 并维持起常驻内存, 若异常退出则再次启动一个
		if xprocess.AlwaysFork2Std(ctx) { // 主进程将阻塞在此
			// 主进程退出
			fmt.Println(pid, " = 管理进程, 退出")
			return
		}

		//子进程可以继续往后执行
		fmt.Println(pid, "= 业务进程")
	}

	// 开始业务逻辑
	start(ctx, cancel)
}

// 开始业务代码
func start(ctx context.Context, cancel func()) {
	name := fmt.Sprintf("业务进程(%d):", os.Getpid())
	idx := 0

	ch := xprocess.NewSignal().SignalAllExit().Done()

	for {
		idx++
		select {
		case <-ctx.Done():
			fmt.Println(name, "ctx.Done(), 退出")
			return
		case <-ch:
			fmt.Println(name, "SignalAllExit(), 退出")
			return
		default:
			log.Println(name, idx)
		}
		time.Sleep(time.Second * 2)
	}
}
