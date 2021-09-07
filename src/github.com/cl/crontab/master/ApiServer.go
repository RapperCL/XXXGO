package master

import (
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
func handlerJobSave(w http.ResponseWriter, r *http.Request) {

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
