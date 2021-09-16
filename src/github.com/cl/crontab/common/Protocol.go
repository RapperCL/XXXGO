package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// 定义定时任务job
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

// 统一返回结构体
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 变化事件
type JobEvent struct {
	EventType int
	Job       *Job
}

// 任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 //任务信息
	Expr     *cronexpr.Expression // cron表达式
	NextTime time.Time            // 下次执行时间
}

//任务执行状态,  应该保存当前任务的执行信息 时间和任务即可
type JobExecuteInfo struct {
	Job      *Job      //任务信息
	PlanTime time.Time // 理论上的调度时间
	RealTime time.Time // 实际的调度时间

	Cancelctx  context.Context    // 任务command的context
	CancelFunc context.CancelFunc // 用于强行取消任务执行
}

// 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // 执行状态
	Output      []byte          //脚本输出
	Err         error
	StartTime   time.Time //启动时间
	EndTime     time.Time //结束时间

}

// 构建应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	// 1 定义一个repsone
	var (
		response Response
	)
	response.Msg = msg
	response.Errno = errno
	response.Data = data

	// 2 序列化json
	if resp, err = json.Marshal(response); err != nil {
		return
	}
	return
}

// 反序列化job
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)
	job = &Job{}
	fmt.Println("要被序列化的输出", string(value))
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

// 任务变化事件
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

// 从目录结构中提取jobName
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// 从目录中提取任务名
func ExtractKillerName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_KILLER_DIR)
}

// 创建执行任务
func BuildJobExecutePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	// 解析job的cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}
	// 生成执行任务
	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return
}

// 构建执行任务
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime, // 计算调度时间
		RealTime: time.Now(),               // 真是调度时间
	}

	return
}
