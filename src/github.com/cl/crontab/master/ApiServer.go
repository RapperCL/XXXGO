package master

import (
	"FoG/src/github.com/cl/crontab/common"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

//日志查看相关
type ApiServer struct {
	httpServer *http.Server
}

var (
	// 单例对象，首字母大写，可以被其他包访问到
	G_apiServer *ApiServer
)

// 保存任务的接口
// POST  job={name  command  cronExpr}
func handlerJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	// 1 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	// 2 获取表单的数据
	postJob = req.PostForm.Get("job")
	// 3 反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	// 4 调用etcd层存放job数据
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	// 5 返回正常应答-- 定义统一响应结构体
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		// 写入响应
		resp.Write(bytes)
	} else {
		goto ERR
	}

	return

ERR:
	fmt.Println(err)
}

// 查询任务的接口
func handlerJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*common.Job
		bytes   []byte
		err     error
	)
	//基于前缀去查询
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	//正常应答-- 将对象序列化为json类型
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		resp.Write(bytes)
	} else {
		goto ERR
	}

	return
ERR:
	fmt.Println("接口获取任务信息错误", err)
	return
}

// 删除服务
func handlerJobDel(resp http.ResponseWriter, req *http.Request) {
	// 基于jobname进行删除
	var (
		err     error
		jobName string
		jobList []*common.Job
		bytes   []byte
	)
	if err = req.ParseForm(); err != nil {
		fmt.Println("删除服务-- 表单解析失败")
		goto ERR
	}
	jobName = req.PostForm.Get("name")
	if jobList, err = G_jobMgr.DelJob(jobName); err != nil {
		goto ERR
	}

	// 获取jobList，然后序列化赋值给表单

	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		resp.Write(bytes)
	} else {
		goto ERR
	}
	return

ERR:
	fmt.Println("删除服务异常", err)
	return
}

// 强杀服务
func handlerJobKill(resp http.ResponseWriter, req *http.Request) {

}

//初始化服务
func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	// 配置路由服务
	mux = http.NewServeMux()
	//  路由,目标方法
	mux.HandleFunc("/job/save", handlerJobSave)
	// 查询服务
	mux.HandleFunc("/job/list", handlerJobList)
	// 删除服务
	mux.HandleFunc("/job/del", handlerJobDel)
	// 任务强杀
	mux.HandleFunc("/job/kill", handlerJobKill)

	// 启动TCP监听 , 并将产生的错误抛出
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个Http服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Second,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Second,
		Handler:      mux,
	}

	// 单例赋值
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动了http服务端
	go httpServer.Serve(listener)

	return
}
