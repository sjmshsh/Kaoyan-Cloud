package count

import (
	"github.com/gin-gonic/gin"
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
	InitRpcCountClient()
	h := New()
	// 这个接口是用于增加计数的，例如我点赞了一次，评论了一次，等等，都统一用这个接口
	r.GET("/count", h.Count)
	// 获取计数信息，也就是有多少个了
	r.GET("/getcount", h.GetCount)
}
