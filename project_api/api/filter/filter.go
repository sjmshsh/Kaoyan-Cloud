package filter

import (
	"context"
	"github.com/gin-gonic/gin"
	pb "github.com/sjmshsh/grpc-gin-admin/project_api/api/filter/protoc"
	"github.com/sjmshsh/grpc-gin-admin/project_common"
	"github.com/sjmshsh/grpc-gin-admin/project_common/errs"
	"net/http"
	"time"
)

type HandlerFilter struct {
}

func New() *HandlerFilter {
	return &HandlerFilter{}
}

func (h *HandlerFilter) filter(ctx *gin.Context) {
	result := &project_common.Result{}
	content := ctx.PostForm("content")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := FilterServiceClient.Filter(c, &pb.ContentMessage{
		Content: content,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(resp.AfterContent))
}
