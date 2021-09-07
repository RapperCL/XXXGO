package main

import (
	"fmt"
	"time"
)

func main() {
	var ch = make(chan bool)

	go func() {
		select {
		case <-ch:
			fmt.Println("2  ")
		}
	}()
	fmt.Println("1 ")
	time.Sleep(2 * time.Second)
	ch <- true
	fmt.Println("3 ")
}
