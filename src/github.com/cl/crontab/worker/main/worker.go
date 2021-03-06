package main

import (
	"FoG/src/github.com/cl/crontab/worker"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string
)

//解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./worker.json", "worker.json")
	flag.Parse()
}

//初始化线程数
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	//初始化命令行
	initArgs()

	//初始化线程
	initEnv()

	//加载配置
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}
	//启动执行器
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}
	// 初始化任务调度器
	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}
	// 初始化任务处理器
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	// 正常退出
	for {
		time.Sleep(worker.G_config.WorkerSleepTime * time.Second)
	}

	return

ERR:
	fmt.Println(err)
}
