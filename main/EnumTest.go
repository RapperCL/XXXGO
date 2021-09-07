package main

import "fmt"

/**
golang 枚举类，要知道在golang中是没有枚举关键字的，不像java，
但是我们可以自定义枚举类
const  name type
const va  int等
*/

type Policy int32

const (
	Policy_one   Policy = 1
	Policy_two   Policy = 2
	Policy_three Policy = 3
)

func foo(policy Policy) {
	fmt.Println("enum value: %v\n", policy)
}

// 对于定于了String()方法的类型，默认输出的时候会调用该方法，实现字符串的打印。例如下面代码
// 于是在方法中执行Println时，就会执行String()方法
func (p Policy) String() string {
	switch p {
	case Policy_one:
		return "one"
	case Policy_two:
		return "two"
	case Policy_three:
		return "three"
	default:
		return "zero"
	}
}

func main() {
	foo(Policy_one)
	fmt.Println()
}
