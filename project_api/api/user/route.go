package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/middleware"
	"github.com/sjmshsh/grpc-gin-admin/project_api/router"
	"log"
)

func init() {
	log.Println("init user router")
	router.Register(&RouterUser{})
}

type RouterUser struct {
}

func (*RouterUser) Router(r *gin.Engine) {
	// 初始化grpc客户端连接
	InitRpcUserClient()
	h := New()
	r.GET("/checkfileMd5", h.CheckFileMd5, middleware.Uv(), middleware.JWTAuth(), middleware.Cors())
	r.GET("/upload", h.UploadFile, middleware.Uv(), middleware.JWTAuth(), middleware.Cors())
	r.GET("/download", h.DownLoadFile, middleware.Uv(), middleware.JWTAuth(), middleware.Cors())
	r.GET("/check", h.CheckIn, middleware.Uv(), middleware.JWTAuth(), middleware.Cors())
	r.GET("/getSign", h.GetSign, middleware.Uv(), middleware.JWTAuth(), middleware.Cors())
	r.GET("/watchUv", h.WatchUv)
	r.GET("/location", h.Location)
	r.POST("/postblog", h.PostBlog)
	r.GET("/watch", h.Watch)

	r.POST("/comment", h.Comment)
	// 关注列表
	r.GET("/list", h.List)
	// 其他各种列表的集合
	r.GET("/olist", h.OList)
	r.GET("/feedlist", h.GetFeedList)
}
