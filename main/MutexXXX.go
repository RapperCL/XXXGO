package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex

	var count = 0
	var wg sync.WaitGroup // 类似一个线程池
	wg.Add(10)            // 定义10 个线程

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done() // 开始执行
			for j := 0; j < 10; j++ {
				mu.Lock() // 获取锁，未获取到的等待
				count++
				mu.Unlock()
			}
		}()
	}
	wg.Wait() // 等待所有的执行完成
	fmt.Println(count)
}
