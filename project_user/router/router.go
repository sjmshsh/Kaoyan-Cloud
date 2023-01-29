package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_common/discovery"
	"github.com/sjmshsh/grpc-gin-admin/project_common/logs"
	"github.com/sjmshsh/grpc-gin-admin/project_user/config"
	user_service_v12 "github.com/sjmshsh/grpc-gin-admin/project_user/pkg/service/user_service_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
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

func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addr, logs.LG)
	resolver.Register(etcdRegister)

	info := discovery.Server{
		Name:    config.C.GC.Name,
		Addr:    config.C.GC.Addr,
		Version: config.C.GC.Version,
		Weight:  config.C.GC.Weight,
	}
	r := discovery.NewRegister(config.C.EtcdConfig.Addr, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}

func RegisterGrpc() *grpc.Server {
	c := gRPCConfig{
		Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			user_service_v12.RegisterUserServiceServer(g, user_service_v12.New())
		},
	}
	s := grpc.NewServer()
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", config.C.GC.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return s
}
