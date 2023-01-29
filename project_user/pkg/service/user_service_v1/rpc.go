package user_service_v1

import (
	"github.com/sjmshsh/grpc-gin-admin/project_api/config"
	"github.com/sjmshsh/grpc-gin-admin/project_common/discovery"
	"github.com/sjmshsh/grpc-gin-admin/project_common/logs"
	"github.com/sjmshsh/grpc-gin-admin/project_user/pkg/service/count_service_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
)

var CountServieClient count_service_v1.CountServiceClient

func InitRpcCountClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addr, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///count", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	CountServieClient = count_service_v1.NewCountServiceClient(conn)
}

