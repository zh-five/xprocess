package xprocess

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const ENV_NAME = "XW_DAEMON_IDX"

// fork 计数, 序号
var forkIdx int = 0

type CmdOption func(*exec.Cmd)

// 保持完全一样的参数, 启动一个子进程
// 返回 cmd != nil 时为父进程, 否则为子进程
// 无参数时,子进程的所有输出(stdout, stderr)默认会被抛弃
func Fork(opts ...CmdOption) (cmd *exec.Cmd, err error) {
	defer func() {
		forkIdx++
	}()

	//子进程, 退出
	if IsForkPassing() {
		return nil, nil
	}

	cmd = &exec.Cmd{
		Path:   os.Args[0],
		Args:   os.Args,
		Stdout: io.Discard,
		Stderr: io.Discard,
		Env:    os.Environ(),
	}

	forkProcAttr(cmd) // 根据平台设置特别属性

	for _, opt := range opts {
		opt(cmd)
	}

	//设置子进程环境变量
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", ENV_NAME, forkIdx+1))

	//启动子进程
	err = cmd.Start()
	if err != nil {
		println(os.Getpid(), "启动子进程失败:", err)
		return
	} else {
		//执行成功
		println(os.Getpid(), ":", "启动子进程成功:", "->", cmd.Process.Pid, "\n ")
	}

	return
}

// 是否为fork时的路过
func IsForkPassing() bool {
	envIdx, err := strconv.Atoi(os.Getenv(ENV_NAME))
	if err != nil {
		envIdx = 0
	}
	println(os.Getpid(), "forkIdx:", forkIdx, "envIdx:", envIdx)

	return forkIdx < envIdx
}

// fork 一个子进程, 子进程的输出(stdout,stderr)写入到日志文件
func Fork2Log(logFile string) (cmd *exec.Cmd, err error) {
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		println(os.Getpid(), ": 打开日志文件错误:", err)
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
// 子进程的所有输出, 写入到日志文件

// 父进程阻塞, 直到 ctx 被取消时, kill 子进程,返回 true
// 子进程不阻塞, 返回 false
func AlwaysFork2Log(ctx context.Context, logFile string) bool {
	return loopFork(ctx, func() (*exec.Cmd, error) {
		return Fork2Log(logFile)
	})
}

// 父进程阻塞, 直到 ctx 被取消时, kill 子进程,返回 true
// 子进程不阻塞, 返回 false
func AlwaysFork2Std(ctx context.Context) bool {
	return loopFork(ctx, func() (*exec.Cmd, error) {
		return Fork2Std()
	})
}

// 父进程阻塞, 直到 ctx 被取消时, kill 子进程,返回 true
// 子进程不阻塞, 返回 false
func AlwaysFork(ctx context.Context) bool {
	return loopFork(ctx, func() (*exec.Cmd, error) {
		return Fork()
	})
}

// 父进程阻塞, 直到 ctx 被取消时, kill 子进程,返回 true
// 子进程不阻塞, 返回 false
func loopFork(ctx context.Context, fork func() (*exec.Cmd, error)) bool {
	for {
		cmd, err := fork()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if cmd == nil { // 子进程
			return false
		}

		// 等待子进程结束
		println(os.Getpid(), "等待子进程结束")
		ch := waitDone(cmd)

		select {
		case <-ctx.Done():
			cmd.Process.Kill() // 杀掉子进程
			return true
		case <-ch: // 子进程退出

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
