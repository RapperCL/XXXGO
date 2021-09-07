package main

import "fmt"

/**
range关键字
*/

func main() {
	nums := []int{1, 2, 3}
	sum := 0
	// 用于数组： 下标索引（不使用时，用下划线代替）,遍历元素 range 被遍历对象
	for _, num := range nums {
		sum += num
	}
	fmt.Println(sum)

	for i, num := range nums {
		if num == 2 {
			fmt.Println("index:", i)
		}
	}

	// 用于map集合,
	kvs := map[string]string{"a": "one", "b": "two"}
	for k, v := range kvs {
		fmt.Printf("%s--->%s\n", k, v)
	}

}
