package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	auth_service_v1 "github.com/sjmshsh/grpc-gin-admin/project_api/api/auth/protoc"
	"github.com/sjmshsh/grpc-gin-admin/project_common"
	"github.com/sjmshsh/grpc-gin-admin/project_common/errs"
	"net/http"
	"time"
)

type HandlerAuth struct {
}

func New() *HandlerAuth {
	return &HandlerAuth{}
}

func (h *HandlerAuth) GetCode(ctx *gin.Context) {
	result := &project_common.Result{}
	phone := ctx.PostForm("phone")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := AuthServiceClient.GetCode(c, &auth_service_v1.GetCodeRequest{Phone: phone})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	// 这个是状态码，不是我的验证码
	ctx.JSON(http.StatusOK, result.Success(resp.Code))
}

func (h *HandlerAuth) Login(ctx *gin.Context) {
	result := &project_common.Result{}
	userName := ctx.PostForm("username")
	password := ctx.PostForm("password")
	useragent := ctx.PostForm("useragent")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := AuthServiceClient.Login(c, &auth_service_v1.LoginRequest{
		Username:  userName,
		UserAgent: useragent,
		Password:  password,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(resp.Token))
}

func (h *HandlerAuth) Register(ctx *gin.Context) {
	result := &project_common.Result{}
	userName := ctx.PostForm("username")
	password := ctx.PostForm("password")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := AuthServiceClient.Register(c, &auth_service_v1.RegisterRequest{
		Username: userName,
		Password: password,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(resp.Token))
}

func (h *HandlerAuth) Phone(ctx *gin.Context) {
	result := &project_common.Result{}
	phone := ctx.PostForm("phone")
	useragent := ctx.PostForm("useragent")
	code := ctx.PostForm("code")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := AuthServiceClient.Phone(c, &auth_service_v1.PhoneRequest{
		Phone:     phone,
		UserAgent: useragent,
		Code:      code,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(resp.Token))
}

func (h *HandlerAuth) Logout(ctx *gin.Context) {
	result := &project_common.Result{}
	userAgent := ctx.Query("useragent")
	token := ctx.Request.Header.Get("token")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := AuthServiceClient.Logout(c, &auth_service_v1.LogoutRequest{
		UserAgent: userAgent,
		Token:     token,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(resp.Code))
}
