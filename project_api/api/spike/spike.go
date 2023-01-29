package spike

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_api/api/spike/common"
	spike_service_v1 "github.com/sjmshsh/grpc-gin-admin/project_api/api/spike/protoc"
	"github.com/sjmshsh/grpc-gin-admin/project_api/api/spike/response"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/util"
	"github.com/sjmshsh/grpc-gin-admin/project_common"
	"github.com/sjmshsh/grpc-gin-admin/project_common/errs"
	"log"
	"net/http"
	"strconv"
	"time"
)

type HandlerSpike struct {
}

func New() *HandlerSpike {
	return &HandlerSpike{}
}

func (h *HandlerSpike) SendRedPack(ctx *gin.Context) {
	token := ctx.Request.Header.Get("token")
	parseToken, err := util.ParseToken(token)
	if err != nil {
		log.Println(err)
	}
	userId := parseToken.Uid
	result := &project_common.Result{}
	Samount := ctx.Query("amount")
	amount, err := strconv.ParseFloat(Samount, 32)
	if err != nil {
		log.Println(err)
	}
	Snumber := ctx.Query("number")
	number, err := strconv.ParseInt(Snumber, 10, 64)
	if err != nil {
		log.Println(err)
	}
	if amount < float64(number)*0.1 {
		ctx.JSON(common.AmountErr, result.Fail(common.AmountErr, "您输入的金额无法完成分配"))
	}
	// c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	resp, err := SpikeServiceClient.SendRedPack(context.Background(), &spike_service_v1.SendRequest{
		Amount: float32(amount),
		Number: int32(number),
		UserID: int64(userId),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, response.RedPackResponse{
		UserId:    1, // 这个获取UserId的过程是网关做的事情
		RedPackId: resp.RedPacketId,
	})
}

func (h *HandlerSpike) RecvRedPack(ctx *gin.Context) {
	result := &project_common.Result{}
	Sid := ctx.Query("id")
	id, err := strconv.ParseInt(Sid, 10, 64)
	if err != nil {
		log.Println(err)
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := SpikeServiceClient.RecvRedPack(c, &spike_service_v1.RecvRequest{
		Id:     id,
		UserId: 1,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(resp.Msg))
}

func (h *HandlerSpike) ListRedPack(ctx *gin.Context) {
	result := &project_common.Result{}
	Sid := ctx.Query("id")
	id, err := strconv.ParseInt(Sid, 10, 64)
	if err != nil {
		log.Println(err)
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := SpikeServiceClient.ListRedPack(c, &spike_service_v1.ListRequest{Id: id})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(resp.List))
}
