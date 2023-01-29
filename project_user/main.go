package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_common"
	"github.com/sjmshsh/grpc-gin-admin/project_user/config"
	"github.com/sjmshsh/grpc-gin-admin/project_user/pkg/dao"
	s "github.com/sjmshsh/grpc-gin-admin/project_user/pkg/service/user_service_v1"
	"github.com/sjmshsh/grpc-gin-admin/project_user/router"
)

func main() {
	r := gin.Default()
	config.InitConfig()
	router.InitRouter(r)
	// grpc服务注册
	gc := router.RegisterGrpc()
	// grpc服务注册到etcd中
	router.RegisterEtcdServer()
	s.InitRpcCountClient()
	s.InitSignConsumer()
	dao.InitRedis()
	dao.InitMysql()
	stop := func() {
		gc.Stop()
	}
	project_common.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
