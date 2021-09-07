package main

import (
	"fmt"
	"reflect"
)

/**
探讨go中的接口
*/
type doInter interface {
	doOne() string
	doTwo()
}

type childone struct {
	a    int
	name string
}

func (*childone) doOne() string {
	fmt.Println("方法一执行")
	return "s"
}

func (*childone) doTwo() {
	fmt.Println("方法二执行")
}

func main() {
	//var doi doInter;
	//child :=  &childone{ }
	//doi = child
	//fmt.Println(doi.doOne())
	//valueType()

	var ct doInter
	if ct == nil {
		fmt.Println("ct为空")
	}
	param := make([]interface{}, 3)
	param[0] = 88
	param[1] = "btman"
	param[2] = 89

	for index, v := range param {
		if _, ok := v.(int); ok {
			fmt.Printf("params[%d]是int类型", index)
		} else if _, ok := v.(string); ok {
			fmt.Printf("params[%d]是string类型", index)
		} else {
			fmt.Printf("params[%d]是未知类型", index)
		}
	}
}

func valueType() {
	// 返回的是值对象
	child := childone{}
	child.a = 2
	child.name = "asdf"
	var ch = &child
	fmt.Printf("%t", reflect.TypeOf(child))
	fmt.Printf("%t", reflect.TypeOf(ch))
}
