package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/zh-five/xprocess"
)

func main() {
	isKill := flag.Bool("k", false, "kill old")
	flag.Parse()

	if *isKill {
		killOld()
	} else {
		exitSelf()
	}
}

func exitSelf() {
	flag := "sfdsdsaf"
	if !xprocess.UniqueCheck(flag) {
		fmt.Println("已经有进程在运行, 退出")
		return
	}

	fmt.Println("run ...")

	time.Sleep(time.Second * 200)
}

func killOld() {
	flag := "sfdsdsaf"
	xprocess.UniqueCheckAndKillOld(flag, func() {
		fmt.Println("有新进程启动, 退出")
		os.Exit(55)
	})

	fmt.Println("run ...")
	select {}
}
