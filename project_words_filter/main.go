package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_common"
	"github.com/sjmshsh/grpc-gin-admin/project_words_filter/config"
	filter_service_v1 "github.com/sjmshsh/grpc-gin-admin/project_words_filter/pkg/service"
	"github.com/sjmshsh/grpc-gin-admin/project_words_filter/router"
)

func main() {
	r := gin.Default()
	config.InitConfig()
	router.InitRouter(r)
	// 把AC自动机在服务启动的时候就直接建立好
	filter_service_v1.InitAcMachine()
	// grpc服务注册
	gc := router.RegisterGrpc()
	// grpc服务注册到etcd中
	router.RegisterEtcdServer()
	stop := func() {
		gc.Stop()
	}
	project_common.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
