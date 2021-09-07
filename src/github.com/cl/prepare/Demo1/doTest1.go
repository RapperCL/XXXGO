package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("1")
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("2")
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("3")
}
