package main

import (
	"fmt"
	"sync"
	"time"
)

/**
golang中的 waitGroup，其实就是类似java中的CountDownLatch
*/

type Counters struct {
	mu    sync.Mutex
	count uint64
}

// 给 Count类添加 Incr方法
func (c *Counters) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}
func (c *Counters) Count() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
func worker(c *Counters, wg *sync.WaitGroup) {
	wg.Done() // 执行一次，减少1
	time.Sleep(time.Second)
	c.Incr()
}

func main() {
	var counter Counters
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go worker(&counter, &wg)
	}
	// 检查点,等待gorouting都执行完毕
	wg.Wait() // 阻塞等待 由10减少为0

	fmt.Println("扣减完成", counter.Count())

}
