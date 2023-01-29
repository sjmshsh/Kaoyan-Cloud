package auth

import (
	auth_service_v1 "github.com/sjmshsh/grpc-gin-admin/project_api/api/auth/protoc"
	"github.com/sjmshsh/grpc-gin-admin/project_api/config"
	"github.com/sjmshsh/grpc-gin-admin/project_common/discovery"
	"github.com/sjmshsh/grpc-gin-admin/project_common/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
)

var AuthServiceClient auth_service_v1.AuthServiceClient

func InitRpcAuthClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addr, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///count", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	AuthServiceClient = auth_service_v1.NewAuthServiceClient(conn)
}
