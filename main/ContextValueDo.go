package main

import (
	"context"
	"fmt"
)

/*
基于context实现上下文传值
*/

func doValue(ctx context.Context) {
	// 取出键值对
	session, ok := ctx.Value("session").(int)
	fmt.Println(ok)
	if !ok {
		fmt.Println("wrong")
		return
	}
	if session != 1 {
		fmt.Println("session 错误")
		return
	}
	// 取出键值对
	traceID := ctx.Value("trace_id").(string)
	fmt.Println("traceId：", traceID, "--session:", session)
}

func main() {
	ctx := context.WithValue(context.Background(), "trace_id", "88888888")
	// context中添加kv键值对
	ctx = context.WithValue(ctx, "session", 1)
	doValue(ctx)
}
