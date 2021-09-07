package main

import (
	"fmt"
	"sync"
)

func tt(wg *sync.WaitGroup) {
	wg.Done()
	fmt.Println("执行扣减")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go tt(&wg)
	}
	fmt.Println("等待")
	wg.Wait()
	fmt.Println("完成")
}
