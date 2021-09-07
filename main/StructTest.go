package main

import "fmt"

/**
结构体
*/

type Tq struct {
	a int64
	b string
}

func main() {
	var tq Tq
	tq = Tq{1, "a"}
	fmt.Println("----", tq.a, tq.b)
	stest(tq)
	fmt.Println("----", tq.a, tq.b)
	szz(&tq)
	fmt.Println("----", tq.a, tq.b)
}

var globalP int64

// 与java中不一样，java中基本类型为值传递，类类型为引用传递
// 在go中，引用传递可以通过指针实现，其他情况默认都是值传递
func stest(tq Tq) {
	tq.a = 2
	tq.b = "b"
	fmt.Println(tq.a, tq.b)
}

func szz(tq *Tq) {
	tq.a = 2
	tq.b = "b"
	fmt.Println(tq.a, tq.b)
}
