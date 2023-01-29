package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Router 接口
type Router interface {
	Router(r *gin.Engine)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	for _, ro := range routers {
		ro.Router(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(server *grpc.Server)
}
