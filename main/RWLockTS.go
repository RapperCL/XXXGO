package main

import (
	"fmt"
	"sync"
)

/**
RWMutex 读写锁，同样没有记录当前是谁持有此锁，所以如果使用不当容易发生死锁现象
例如，最常见的不可复制
 1 RWMutex 读写锁，是写优先，什么是锁优先？  当前有读锁持有锁时，写锁也会进行等待，只是此时后面来的读锁要等待写锁。
   这就是写锁的优先级体现（此时会不会有其他读锁在等待呢？  不会，因为读读不互斥)
*/

func RWT(mutex sync.RWMutex) {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println("加锁")
}

func main() {
	var mutex sync.RWMutex
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println("加锁")
	//RWT(mutex)
}
