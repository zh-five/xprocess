package fork

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const ENV_NAME = "XW_DAEMON_IDX"

// 运行时调用background的次数
var runIdx int = 0

type CmdOption func(*exec.Cmd)

// 保持完全一样的参数, 启动一个子进程
// 返回 cmd != nil 时为父进程, 否则为子进程
// 无参数时,子进程的所有输出(stdout, stderr)默认会被抛弃
func Fork(opts ...CmdOption) (cmd *exec.Cmd, err error) {
	//判断子进程还是父进程
	runIdx++
	envIdx, err := strconv.Atoi(os.Getenv(ENV_NAME))
	if err != nil {
		envIdx = 0
	}
	if runIdx <= envIdx { //子进程, 退出
		return nil, nil
	}

	cmd = &exec.Cmd{
		Path:        os.Args[0],
		Args:        os.Args,
		SysProcAttr: NewSysProcAttr(),
		Stdout:      io.Discard,
		Stderr:      io.Discard,
	}

	for _, opt := range opts {
		opt(cmd)
	}

	//设置子进程环境变量
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", ENV_NAME, runIdx))

	//启动子进程
	err = cmd.Start()
	if err != nil {
		log.Println(os.Getpid(), "启动子进程失败:", err)
		return
	} else {
		//执行成功
		log.Println(os.Getpid(), ":", "启动子进程成功:", "->", cmd.Process.Pid, "\n ")
	}

	return
}

// fork 一个子进程, 子进程的输出(stdout,stderr)写入到日志文件
func Fork2Log(logFile string) (cmd *exec.Cmd, err error) {
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println(os.Getpid(), ": 打开日志文件错误:", err)
		return
	}

	cmd, err = Fork(func(c *exec.Cmd) {
		c.Stderr = f
		c.Stdout = f
	})

	return
}

// fork 一个子进程, 子进程的输出(stdout,stderr)重定向到父进程的输出
func Fork2Std() (cmd *exec.Cmd, err error) {
	cmd, err = Fork(func(c *exec.Cmd) {
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
	})

	return
}

// fork子进程, 若退出总是再fork一个. 保证总有一个子进程在运行
// 父进程会阻塞在本方法内部一直循环, 直到 ctx 取消
// 子进程不阻塞,直接返回

// ctx 取消后, 父进程会杀死正在运行的子进程, 并结束循环
// 子进程的所有输出, 写入到日志文件
func AlwaysFork2Log(ctx context.Context, logFile string) {
	loopFork(ctx, func() (*exec.Cmd, error) {
		return Fork2Log(logFile)
	})
}

func AlwaysFork2Std(ctx context.Context) {
	loopFork(ctx, func() (*exec.Cmd, error) {
		return Fork2Std()
	})
}

func AlwaysFork(ctx context.Context) {
	loopFork(ctx, func() (*exec.Cmd, error) {
		return Fork()
	})
}

func loopFork(ctx context.Context, fn func() (*exec.Cmd, error)) {
	for {
		cmd, err := fn()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if cmd == nil { // 子进程
			return
		}

		ch := waitDone(cmd)

		select {
		case <-ctx.Done():
			cmd.Process.Kill() // 杀掉子进程
			return
		case <-ch:

		}
	}
}

func waitDone(cmd *exec.Cmd) chan struct{} {
	ch := make(chan struct{})
	go func(c *exec.Cmd) {
		c.Wait()
		close(ch)
	}(cmd)

	return ch
}
