package main

import "fmt"

func zzT1() {

}

type test struct {
	c int32
}

func main() {
	a := 10
	fmt.Println("1 a的内存地址：", &a)
	// 获取到a的内存地址，即b作为了a的指针
	b := &a
	fmt.Println("2 指针b的值：", b)
	p := *b
	fmt.Println("3 指针b存放的内存地址对应的值：", p)
	// * 用于取出指针存放的内存地址对应的值，所以不能用在基本类型中，那普通类型呢？ 也不行，仅仅用于取出指针对应的内存地址的值。
	//var t test
	//c :=*&t
	//fmt.Println(c)
	d := 10
	testM(&d)
}

/**
通过这个例子，可以发现go中的指针与c不同，c中的指针是直接以*变量名命名，而go中* 是取指针对应的内存地址的值， & 才是声明指针（
即 将当前变量的内存地址赋值给变量，此变量作为指针）
*/
// 基于指针来修改对应变量的值
func testM(a *int) {
	fmt.Println("指针值存放的内存地址", a)

	*a = 3
	fmt.Println("指针对应的值", *a)
	// 修改的是指针存放的内存地址的值，此指针指向的内存地址不会变化
	fmt.Println("指针值存放的内存地址", a)
}
