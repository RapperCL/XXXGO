package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	c := sync.NewCond(&sync.Mutex{})
	var ready int

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)

			c.L.Lock()
			ready++
			c.L.Unlock()

			log.Println("运动员%d 已准备就绪\n", i)

			// 广播唤醒所有等待 -- notifyAll
			c.Broadcast()
		}(i)
	}

	// 重点1 提前加锁，使用wait要提前加锁，因为在go中wait等待会将当前goroutin加入到等待队列中，然后释放锁
	//所以需要提前获取锁.(释放之后，等待会sigln 或broadCast 唤醒，重新获取到锁， 自此一个wait方法才走完。
	// 所以说wait方法中：先释放锁，然后加入等待队列，等待获取锁。
	c.L.Lock()
	t := 0
	// 重点 2  需要循环判断
	for ready != 10 {
		c.Wait()
		t += 1
		log.Println("裁判员被唤醒:%d 次", t)
	}

	c.L.Unlock()

	log.Println("所有人都准备就绪了")
}
