package main

import (
	"fmt"
	"github.com/petermattis/goid"
	"sync"
)

func te1(wg *sync.WaitGroup) {
	gid := goid.Get()
	defer wg.Done()
	fmt.Println("1111:  ", gid)
}

func te2(wg *sync.WaitGroup) {
	gid := goid.Get()
	defer wg.Done()
	fmt.Println("2222:  ", gid)
}

// 看了文章说，每个func都会对应一个goroutine，机遇goid来验证一下
// 需要在func前面加上 go，此时会主动创建一个协程
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	gid := goid.Get()
	fmt.Println("-----", gid)
	go te1(&wg)
	fmt.Println("-----", gid)
	go te2(&wg)
	fmt.Println("-----", gid)
	wg.Wait()
}
