package spike

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/middleware"
	"github.com/sjmshsh/grpc-gin-admin/project_api/router"
	"log"
)

func init() {
	log.Println("init user router")
	router.Register(&RouterSpike{})
}

type RouterSpike struct {
}

func (*RouterSpike) Router(r *gin.Engine) {
	// 初始化grpc客户端连接
	InitRpcUserClient()
	h := New()
	r.GET("/sendRedPack", h.SendRedPack, middleware.JWTAuth(), middleware.Cors())
	r.GET("/recvRedPack", h.RecvRedPack, middleware.JWTAuth(), middleware.Cors())
	r.GET("/listRedPack", h.ListRedPack, middleware.JWTAuth(), middleware.Cors())
}
