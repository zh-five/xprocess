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

func main() {
	isKill := flag.Bool("k", false, "kill 现在的进程, 有此参数时, 忽略 -d 参数")
	isDaemon := flag.Bool("d", false, "后台运行")
	flag.Parse()
	//ctx, cancel := context.WithCancel(context.Background())

	flagStr := "gv_cdp"
	logFile := filepath.Join(os.TempDir(), fmt.Sprintf("gv_cdp_%s.log", time.Now().Format("20060102_150405")))

	ctx, cancel, isExit := xprocess.Daemon(flagStr, logFile, *isKill, *isDaemon)
	if isExit {
		return
	}

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
			cancel()
			return
		default:
			log.Println(name, idx)
		}
		time.Sleep(time.Second * 2)
	}
}
