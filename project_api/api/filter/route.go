package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/middleware"
	"github.com/sjmshsh/grpc-gin-admin/project_api/router"
	"log"
)

func init() {
	log.Println("init user router")
	router.Register(&RouterFilter{})
}

type RouterFilter struct {
}

func (*RouterFilter) Router(r *gin.Engine) {
	// 初始化grpc客户端连接
	InitRpcUserClient()
	h := New()
	r.POST("/filter", h.filter, middleware.JWTAuth(), middleware.Cors())
}
