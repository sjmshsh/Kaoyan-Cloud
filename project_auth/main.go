package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_auth/config"
	"github.com/sjmshsh/grpc-gin-admin/project_auth/pkg/dao"
	"github.com/sjmshsh/grpc-gin-admin/project_auth/router"
	"github.com/sjmshsh/grpc-gin-admin/project_common"
)

func main() {
	r := gin.Default()
	config.InitConfig()
	router.InitRouter(r)
	// grpc服务注册
	gc := router.RegisterGrpc()
	// grpc服务注册到etcd中
	router.RegisterEtcdServer()
	// 把MySQL连接
	dao.InitMysql()
	// 把redis连接
	dao.InitRedis()
	stop := func() {
		gc.Stop()
	}
	project_common.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
