package main

import (
	"flag"
	"fmt"
	"github.com/cl/crontab/master"
	"runtime"
)

var (
	confFile string
)

// 通过命令行传入文件路径
func initArgs() {
	// ./master.json  设置配置文件解析之后的接收方  默认为当前目录下的master.json文件
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	// 将命令行的参数解析出来并赋值给confFile对象
	flag.Parse()
}

// 环境变量-- 1 设置最大线程数和当前cpu核数一样多
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	// 初始化命令行参数
	initArgs()

	// 初始化线程
	initEnv()

	// 加载配置
	if err = master.InitConfig(""); err != nil {
		goto ERR
	}

	// 启动ApiServer服务  基于包进行调用！
	if err = master.InitApiServer(); err != nil {
		fmt.Println("初始化http服务失败!", err)
		goto ERR
	}

	return
ERR:
	fmt.Println(err)
}
