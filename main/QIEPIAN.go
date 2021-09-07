package main

import "fmt"

/**
切片： 对于数组的抽象，为了弥补数组长度不能动态变化的缺陷，在某些程度上类似于java的动态数组，但是又有区别。

*/

func main() {
	var qp []int
	//println(qp,len(qp),cap(qp))
	fmt.Printf("%v,%d,%d \n", qp, len(qp), cap(qp)) // 空

	var te []int = []int{1, 2, 3, 4, 5}
	fmt.Printf("%v,%d,%d \n", te, len(te), cap(te))

	//截取功能  体现了切片特点
	fmt.Printf("%v,%d,%d \n", te[:3], len(te[:3]), cap(te[:2]))

	var tr []int = te[1:2]
	fmt.Printf("%v,%d,%d \n", tr, len(tr), cap(tr))

	//利用make创建切片
	tt := make([]int, 3, 5)
	fmt.Printf("%v,%d,%d \n", tt, len(tt), cap(tt))

	//切片扩充，类似于java中动态字符串
	tt = append(tt, 1, 2, 3, 4)
	fmt.Printf("%v,%d,%d \n", tt, len(tt), cap(tt))

	//切片拷贝
	ty := make([]int, len(tt), cap(tt)-1)
	copy(ty, tt)
	fmt.Printf("%v,%d,%d \n", ty, len(ty), cap(ty))

}
