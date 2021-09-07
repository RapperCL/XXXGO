package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/**
基于token来实现 Mutex的可重入版本
*/

type MutexTk struct {
	sync.Mutex
	token int64
	count int32
}

// 利用token存放当前goroutine下生成的唯一标识
func (m *MutexTk) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token {
		m.count++
		return
	}

	//非同一个，于是就要进行加锁
	m.Mutex.Lock()
	//获取到锁之后
	m.count++
	//存储起来，供后面进行判断
	atomic.StoreInt64(&m.token, token)

}

func (m *MutexTk) Unlock(token int64) {
	// 首先判断是否持有锁的，不是就抛出异常
	if atomic.LoadInt64(&m.token) != token {
		panic(fmt.Sprint("unlock error: %d, token: %d", m.token, token))
	}

	// 是同一个就减一之后，判断count是否为零了
	m.count--

	if m.count != 0 {
		return
	}
	// 减少为0了，就要真正解锁了
	atomic.StoreInt64(&m.token, -1)
	m.Mutex.Unlock()
}

func main() {
	var token int64
	var token2 int64
	token = 10
	token2 = 20

	fmt.Println("加锁")
	var mt MutexTk
	mt.Lock(token)
	go func() {
		mt.Lock(token2)
		fmt.Println("加锁成功")
		mt.Unlock(token2)
	}()
	defer mt.Unlock(token)
	fmt.Println("释放成功")
	time.Sleep(10000)
}
