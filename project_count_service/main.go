package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_common"
	"github.com/sjmshsh/grpc-gin-admin/project_count_service/config"
	"github.com/sjmshsh/grpc-gin-admin/project_count_service/pkg/dao"
	count_service_v1 "github.com/sjmshsh/grpc-gin-admin/project_count_service/pkg/service"
	"github.com/sjmshsh/grpc-gin-admin/project_count_service/router"
)

func main() {
	r := gin.Default()
	config.InitConfig()
	router.InitRouter(r)
	// grpc服务注册
	gc := router.RegisterGrpc()
	// grpc服务注册到etcd中
	router.RegisterEtcdServer()
	dao.InitMysql()
	dao.InitRedis()
	count_service_v1.InitFuzzyScheduler()
	count_service_v1.ConsumeDeleteFailed()
	stop := func() {
		gc.Stop()
	}
	project_common.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
