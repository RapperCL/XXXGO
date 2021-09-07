package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

func test1() {
	flag := true
	fmt.Println(flag)
	var ss = 3 | 6
	fmt.Println(ss)
	// 多声明
	var n1, n2, n3 = 1, 2, 3
	fmt.Println(n1, n2, n3)
	n4, n5, n6 := 4, 5, 6
	fmt.Println(n4, n5, n6)

	var age int
	age = 10
	var ul = "No:%s,age:%d"
	fmt.Println(age)
	fmt.Println("Hello,world!")
	var tag = fmt.Sprintf(ul, "4557", age)
	fmt.Println(tag)
	test2()
}

// 验证欧拉公式 e的iπ次方+1 =0
func test2() {
	c := 3 + 4i
	fmt.Println(cmplx.Abs(c))
	fmt.Println(cmplx.Pow(math.E, 1i*math.Pi) + 1)
}

var ch int

func main() {
	test1()
	//a,b :=swap("1","2");
	//cc(2);
	//fmt.Println(a,b,"-----",ch+1)
	arr()
	zhizhen()
	a, b := 0, 0
	fmt.Println(&a, "----", &b)
	zhiT(a, b)
	fmt.Println(&a, "----", &b)
	yinT(&a, &b)
	fmt.Println(&a, "----", &b)
	globalP = 10

}

func cc(x int) (ch int) {
	return 4 * x
}

func swap(x, y string) (string, string) {
	return y, x
}

func arr() {
	var n [10]int

	for i := 0; i < 10; i++ {
		n[i] = i + 100
		fmt.Println(n[i])
	}
}

func zhizhen() {
	a := 10
	fmt.Println(&a)
}

func zhiT(a int, b int) {
	a = 3
	b = 4
	fmt.Println("zhiT", &a, "---", &b)
}

// 引用传递，此时修改的是对应的内存地址的值
// 值传递时，会新建一个内存地址来进行存放，并不会修改原来的内存地址
func yinT(a *int, b *int) {
	fmt.Println("yinT1:", &a, "----", &b)
	*a = 3
	*b = 4
	fmt.Println("yinT2:", &a, "----", &b)
}
