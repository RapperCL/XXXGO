package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	output []byte
	err    error
}

func main() {
	// 基于context实现协程通信
	var (
		ctx        context.Context
		cancel     context.CancelFunc
		cmd        *exec.Cmd
		resultChan chan *result // 下面利用channel实现协程之间通信
		res        *result
	)
	ctx, cancel = context.WithCancel(context.TODO())

	// channel-- 结果队列
	resultChan = make(chan *result, 1000)

	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, "E:\\Git\\bin\\bash.exe", "-c", "sleep 10;echo hello;")
		// 基于select函数，去监听ctx.done()  done()会在被取消或关闭时返回
		// 于是我们可以在主函数中，执行cancel函数，此时select就会监听到ctx.done()
		// kill 杀死子进程
		output, err = cmd.CombinedOutput()
		// 将子协程的结果输出到channel中
		resultChan <- &result{
			err:    err,
			output: output,
		}
	}()

	time.Sleep(1 * time.Second)
	// 1s之后中断掉
	// 执行取消函数，变量接收取消函数
	cancel()

	// 在main协程中，等待子协程的退出，并打印任务执行结果
	res = <-resultChan

	fmt.Println(string(res.output), "错误信息", res.err)
}
