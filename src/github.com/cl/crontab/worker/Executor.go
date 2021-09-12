package worker

import (
	"FoG/src/github.com/cl/crontab/common"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Executor struct {

}
var(
	G_executor *Executor
)

func InitExecutor()(err error){
	G_executor =&Executor{}
	return
}
//1 执行任务
func(executor *Executor) ExecuteJob(info *common.JobExecuteInfo){
	go func() {
		var(
			result *common.JobExecuteResult
			err error
			cmd *exec.Cmd
			output []byte
			jobLock *JobLock
		)
		//2 任务结果封装
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output: make([]byte,0),
		}
        // 2-1 初始化分布式锁
        jobLock = G_jobMgr.CreateJobLock(info.Job.Name)

		//3 记录任务开始时间
		result.StartTime = time.Now()

		//2-2 上锁
		err = jobLock.TryLock()
		// 2-3 方法退出时，即任务执行完成时，释放锁
		defer jobLock.Unlock()

		if err != nil{
			result.Err = err
			result.EndTime = time.Now()
			fmt.Println(info.Job.Name,"获取锁失败...")
		}else{
			fmt.Println(info.Job.Name,"获取锁成功...")
			// 获取锁成功之后，执行任务
			//4 执行shell命令
			cmd = exec.CommandContext(context.TODO(),"E:\\Git\\bin\\bash.exe","-c",info.Job.Command)

			//5 执行并获取输出
			if output, err = cmd.CombinedOutput(); err != nil{
				fmt.Println(info.Job.Name,"任务执行失败...","失败原因:",err)
			}else{
				fmt.Println(info.Job.Name,"任务执行成功...")
			}

			//6 记录任务结束时间
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}

		//7 任务执行完成或失败之后，将任务结果返回给scheduler-- 从execute表中移除执行完成的任务
        G_scheduler.PushJobResult(result)
	}()
}