package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zh-five/xprocess/sig"
)

func main() {
	ctx, _ := sig.NewSignal().SignalAllExit().WithCancel(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("end")
			os.Exit(1)
		}
	}()

	select {}
}
