package main

import (
	"fmt"
	"os"

	"github.com/zh-five/xprocess"
)

func main() {

	go func() {
		<-xprocess.NewSignal().SignalAllExit().Done()
		fmt.Println("end") // ctrl +c 后会输出 end
		os.Exit(1)
	}()

	select {}
}
