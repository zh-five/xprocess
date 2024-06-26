package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zh-five/xprocess"
)

func main() {
	ctx, _ := xprocess.NewSignal().SignalAllExit().WithCancel(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("end")
			os.Exit(1)
		}
	}()

	select {}
}
