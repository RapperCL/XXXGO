package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time //expr.Next(now)获取下次执行时间
}

func main() {
	// 1个调度协程来定时检查所有的cron任务，过期之后就执行

	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob
	)
	scheduleTable = make(map[string]*CronJob)
	now = time.Now()
	//定义2个cronJob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now), // 下次执行时间
	}
	//将任务1放入到调度表中
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now), // 下次执行时间
	}
	//将任务2放入到调度表中
	scheduleTable["job2"] = cronJob

	//此时调度表中有两个任务, 启动一个调度协程
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)
		for {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				// 判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协程，去执行这个任务
					go func(jobName string) {
						fmt.Println("执行任务：", jobName)
					}(jobName)

					// 计算下一次执行时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println("下次执行时间:", cronJob.nextTime)
				}
			}

			select {
			// 过期之后，会往channe中发送数据，可读
			case <-time.NewTimer(100 * time.Millisecond).C:
			}
		}
	}()

	time.Sleep(20 * time.Second)
}
