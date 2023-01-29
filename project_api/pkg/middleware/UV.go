package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/common"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/dao"
	"log"
	"net/http"
)

// Uv 用户UV统计的中间件, 使用HyperLogLog
func Uv() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先获取前端HTTP请求头里面发过来的IP地址，我们通过IP地址来进行UV统计
		ip := c.GetHeader("ip")
		if ip == "" {
			// 如果是空的话，说明前端没有传过来信息，这里就不进行统计了
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "请求没有携带用户的相关信息，无法进行统计",
			})
			// 在被调用的函数中阻止后续中间件的执行
			c.Abort()
			return
		}
		// 不为空，进行UV统计
		err := dao.Rdb.PFAdd(dao.RCtx, common.WebUv, ip).Err()
		if err != nil {
			log.Println(err)
		}
		return
	}
}
