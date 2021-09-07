package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 1 建议在函数头部这样集中管理变量
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)

	cmd = exec.Command("E:\\Git\\bin\\bash.exe", "-c", "sleep 5;ls -l")

	// 2 推荐这种方式写法
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))
}
