package middleware

import (
	"github.com/gin-gonic/gin"
	dao2 "github.com/sjmshsh/grpc-gin-admin/project_api/pkg/dao"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/util"
	"net/http"
	"strconv"
	"time"
)

// JWTAuth 定义一个JWTAuth中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "请求未携带token，无权访问",
			})
			// 在被调用的函数中阻止后续中间件的执行
			c.Abort()
			return
		}
		// 解析token中包含的相关信息(有效载荷)
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "token解析失败",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		// 判断该token是不是最新的token(从redis里面查)
		ua := c.GetHeader("User-Agent")
		val, err := dao2.Rdb.HGet(dao2.RCtx, strconv.Itoa(int(claims.Uid)), ua).Result()
		if err != nil {
			// 说明该token是其他User-Agent的token（你如说电脑端的token，当然不能用来登录手机端）
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "token所属的User-Agent不匹配",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		if token != val {
			// 请求携带的token与redis中存储的token不一致，说明是旧的token
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "token失效!",
			})
			c.Abort()
			return
		}
		// 处理过期token
		if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "token已经过期了", // token过期就得重新登录的
			})
			c.Abort()
			return
		}
	}
}
