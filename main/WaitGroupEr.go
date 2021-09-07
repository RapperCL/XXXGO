package main

import (
	"fmt"
	"sync"
)

type TestStruct struct {
	Wait sync.WaitGroup
}

func main() {
	w := sync.WaitGroup{}
	w.Add(1)
	t := &TestStruct{
		Wait: w, // 赋初值,并将WaitGroup做了复制
	}
	t.Wait.Done()
	fmt.Println("Finished")
}
