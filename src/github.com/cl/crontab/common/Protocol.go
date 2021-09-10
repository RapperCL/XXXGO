package common

import (
	"encoding/json"
	"strings"
)

// 定义定时任务job
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

// 统一返回结构体
type Response struct {
	Errno int         `json:"error"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 变化事件
type JobEvent struct {
	EventType int
	Job       *Job
}

// 构建应答方法
func BuildResponse(erron int, msg string, data interface{}) (resp []byte, err error) {
	// 1 定义一个repsone
	var (
		response Response
	)
	response.Msg = msg
	response.Errno = erron
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
