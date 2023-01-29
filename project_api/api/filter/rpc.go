package filter

import (
	pb "github.com/sjmshsh/grpc-gin-admin/project_api/api/filter/protoc"
	"github.com/sjmshsh/grpc-gin-admin/project_api/config"
	"github.com/sjmshsh/grpc-gin-admin/project_common/discovery"
	"github.com/sjmshsh/grpc-gin-admin/project_common/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
)

var FilterServiceClient pb.FilterServiceClient

func InitRpcUserClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addr, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///filter", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	FilterServiceClient = pb.NewFilterServiceClient(conn)
}
