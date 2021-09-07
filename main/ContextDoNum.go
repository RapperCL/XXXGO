package main

import (
	"context"
	"fmt"
	"time"
)

/**
利用context+channel实现个数控制，当消耗了指定的数目时，就停止子goroutine的实现

如果单单以channel，可以实现吗？
只能控制当前的goroutine，只能在一个函数内，除非将ch定义为全局变量，否则就不能在多个函数中传播，
不像context，可以在多个上下文中传输
*/

func onlyChan(ct context.Context, i int) {
	for {
		select {
		case <-ct.Done():
			fmt.Printf("子执行了i%d个，就该停止了", i)
			return
		default:
			i++
			fmt.Printf("子执行了i%d个，执行中", i)
		}
	}
}

func main() {

	ct, cancel := context.WithCancel(context.Background())

	//	 ch <- true
	go onlyChan(ct, 1)
	// 父类控制子类,
	i := 0
	for {
		i++
		time.Sleep(1 * time.Second)
		if i > 10 {

			fmt.Println("父类执行超过了10个，该让子类停止了")
			cancel()
			return
		}
	}

}
