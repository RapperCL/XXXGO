package main

import "fmt"

func defer_call() {
	//defer
	func() {
		fmt.Println("打印前")
	}()

	defer func() {
		fmt.Println("打印中")
	}()

	defer func() {
		fmt.Println("打印后")
	}()

	//panic("触发异常")

}

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func main() {

	println(DeferFunc1(1))
	println(DeferFunc2(1))
	println(DeferFunc3(1))
	println(DeferFunc4(2))

	//a :=1
	//b :=2
	//defer calc("1",a,calc("10",a,b))
	//// 10  1 2 3
	////  1  1 3 4
	//a=0
	//defer calc("2",a,calc("20",a,b))
	//// 20 0 2 2
	//// 2 0 2 2
	//b=1

	//	defer_call();
}

// i = 1, 返回了t
func DeferFunc1(i int) (t int) {
	t = i
	// 入栈时，t = i，在return之后执行，返回4
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc2(i int) int {
	t := i
	defer func() { // 先返回了t，且t不是固定的返回类型 int t,
		// 所以此处为1, 如果这里是  以作用域角度来理解，t的作用域为函数内，当
		// return t执行之后，t的作用域就没有了，即使t+=3 ，但是不会返回，最终返回的是t=1
		t += 3
		println("f2", t)
	}()
	return t
}

func DeferFunc3(i int) (t int) {
	// 此时t未被赋值，需要等返回之后，进行计算 ，return 2== 赋值t=2
	// t的作用域为整个函数体，return 返回时，t=2，defer 计算出 t+1=3 ，最终返回 3
	defer func() {
		t += i
	}()
	return 2
}

func DeferFunc4(i int) int {
	// 此时x的默认值为0，于是输出0
	println("one:", i)

	// 在return之后执行，返回了4，代表x = 4
	defer func() {
		i *= i
	}()
	return 4
}
