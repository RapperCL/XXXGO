package main

import (
	"sync"
)

type Counter struct {
	sync.Mutex
	Count int
}

func main() {
	var c Counter
	c.Lock()
	defer c.Unlock()
	c.Count++
	//此时会一直等待，因为是不可重入的，但不往下走，又无法解锁，，，于是抛出deadLock异常
	c.Lock()
	c.Unlock() //此时无

	//foo(c)
}

/**
在java中，的reentrantLock锁，是不会报错的，包括synchronize，因为他们都是可重入的，
他们之所以能够可重入，就拿reentrantLock来说，毫无疑问是因为他通过线程id记录了当前是谁获取了这把锁。
ThreadOwen
mutex不能的原因就是因为没有记录，那我们是否可以人为进行改造呢？
*/
//func foo(c ReMutex){
//	// 前面已经加锁了，这里二次加锁，重复加锁，加锁失败，且无法解锁。
//	c.Lock()
//	defer c.Unlock()
//	fmt.Println("in foo")
//}
