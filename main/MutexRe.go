package main

import (
	"fmt"
	"github.com/petermattis/goid"
	"sync"
	"sync/atomic"
)

/**
基于Mutex并结合reentrantLock的可重入特性，实现可重入锁
*/

type ReMutex struct {
	sync.Mutex
	owner int64 //go routin id 线程id
	count int32 // 重入的次数
}

func (m *ReMutex) Lock() {
	gid := goid.Get() // 获取当前的gid
	if atomic.LoadInt64(&m.owner) == gid {
		m.count++
		return
	}
	// 非拥有锁的goroutin
	m.Mutex.Lock()
	// 走到这里，都是第一次获取成功的，gid赋值给owner
	atomic.StoreInt64(&m.owner, gid)
	m.count = 1
}

// 解锁，首先判断是否拥有锁的goroutin，持有锁的才能释放
func (m *ReMutex) Unlock() {
	gid := goid.Get()
	//非持有锁的gor，不让释放，并抛出异常
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d", m.owner, gid))
	}
	//走到这里代表是拥有锁的goroutin
	m.count--

	if m.count != 0 {
		return
	}
	// 减到0时，释放
	atomic.StoreInt64(&m.owner, -1)
	m.Unlock()
}

//
