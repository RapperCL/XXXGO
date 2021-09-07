package common

import "encoding/json"

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
