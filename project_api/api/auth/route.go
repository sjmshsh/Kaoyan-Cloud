package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/middleware"
	"github.com/sjmshsh/grpc-gin-admin/project_api/router"
	"log"
)

func init() {
	log.Println("init auth router")
	router.Register(&RouterAuth{})
}

type RouterAuth struct {
}

func (*RouterAuth) Router(r *gin.Engine) {
	// 初始化grpc客户端连接
	InitRpcAuthClient()
	h := New()
	r.POST("/login", h.Login, middleware.Cors())
	r.POST("/register", h.Register, middleware.Cors())
	r.POST("/phone", h.Phone, middleware.Cors())
	r.GET("/getcode", h.GetCode, middleware.Cors())
	r.GET("/logout", h.Logout, middleware.Cors())
}
