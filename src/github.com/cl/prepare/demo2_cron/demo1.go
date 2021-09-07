package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)

	//每隔5分钟执行一次  分钟  小时 天 月 星期  实现原理，枚举 0 5 10
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	// 当前时间
	now = time.Now()

	// 下次执行时间
	nextTime = expr.Next(now)

	// 等待定时器超时
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了:", nextTime)
	})

	time.Sleep(8 * time.Second)
}
