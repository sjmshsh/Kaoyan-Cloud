package count

import (
	"context"
	"github.com/gin-gonic/gin"
	count_service_v1 "github.com/sjmshsh/grpc-gin-admin/project_api/api/count/protoc"
	"log"
	"net/http"
	"strconv"
)

type HandlerCount struct {
}

func New() *HandlerCount {
	return &HandlerCount{}
}

func (h *HandlerCount) Count(ctx *gin.Context) {
	ty := ctx.Query("type")
	id := ctx.Query("id")
	s := ctx.Query("symbol")
	symbol, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println(err)
	}
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
	}
	// t是0代表是模糊计数
	// 1代表是精准计数
	typ, err := strconv.ParseInt(ty, 10, 64)
	if err != nil {
		log.Println(err)
	}
	req := &count_service_v1.CountRequest{
		Type:   typ,
		Id:     i,
		Symbol: symbol,
	}
	resp, err := CountServieClient.Count(context.Background(), req)
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, resp.Data)
}

func (h *HandlerCount) GetCount(ctx *gin.Context) {
	ty := ctx.Query("type")
	id := ctx.Query("id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
	}
	typ, err := strconv.ParseInt(ty, 10, 64)
	if err != nil {
		log.Println(err)
	}
	req := &count_service_v1.CountRequest{
		Type: typ,
		Id:   i,
	}
	resp, err := CountServieClient.GetCount(context.Background(), req)
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, resp.Data)
}
