package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	//并发场景利用channel,实现多个goroutine之间的通信，这些goroutine之间是平级的，可以通过channel通道来实现数据交换
	//chanDo();
	doContext()
}

func chanDo() {
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("3 监控退出，停止了....")
				return
			default:
				fmt.Println("1 goroutine监控中...")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("2 通知监控停止")
	stop <- true
	time.Sleep(5 * time.Second)
}

func doContext() {
	// 传递一个父context作为参数，返回子context，以及一个取消函数 cancel
	ctx, cancel := context.WithCancel(context.Background())

	go watch(ctx, "【监控1】")
	//go watch(ctx,"【监控2】")
	//go watch(ctx,"【监控3】")
	time.Sleep(10 * time.Second)
	fmt.Println("2 通知监控停止")
	// 执行取消时，会发送数据到channel吗？
	cancel()
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context, name string) {
	for {
		select { // select函数： 类似switch case一样
		// 从channel中读取  -- context 的done方法会返回一个channel
		case <-ctx.Done():
			fmt.Println(name, "3 监控退出了...")
			return
		default:
			fmt.Println(name, "1 goroutine监控中...")
			time.Sleep(2 * time.Second)
		}
	}
}
