package auth_service_v1

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/sjmshsh/grpc-gin-admin/project_auth/pkg/common"
	"github.com/sjmshsh/grpc-gin-admin/project_auth/pkg/dao"
	"github.com/sjmshsh/grpc-gin-admin/project_auth/pkg/model"
	"github.com/sjmshsh/grpc-gin-admin/project_auth/pkg/util"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthService struct {
	UnimplementedAuthServiceServer
}

func New() *AuthService {
	return &AuthService{}
}

func (AuthService) GetCode(ctx context.Context, req *GetCodeRequest) (*GetCodeResponse, error) {
	redisCode, err := dao.Rdb.Get(dao.RCtx, common.AuthCode+req.Phone).Result()
	if !errors.Is(err, redis.Nil) {
		// 如果没有这个错误的话，就说明我的redis里面已经有验证码了
		// 我做出相应的处理，避免接口被刷
		split := strings.Split(redisCode, "_")
		s := split[1]
		redistime, _ := strconv.ParseInt(s, 10, 64)
		if time.Now().UnixNano()-redistime < 60*1000 {
			// 60s
			return &GetCodeResponse{Code: common.CodeRepeat}, nil
		}
	}
	// 我的redis里面没有验证码，那么此时就要往redis里面加入验证码了
	phone := req.Phone
	// 生成验证码
	code := util.Code()
	// 把验证码保存到redis里面去
	// 设计5分钟的过期时间
	dao.Rdb.Set(dao.RCtx, common.AuthCode+phone, code, time.Minute*5)
	// TODO 整合阿里云的短信服务
	return &GetCodeResponse{
		Code: http.StatusOK,
	}, nil
}

func (AuthService) Login(ctx context.Context, req *LoginRequest) (*Response, error) {
	var user model.User
	result := dao.MysqlDB.Where(&model.User{UserName: req.Username}).First(&user)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果数据库中没有找到记录
			return &Response{
				Code: common.UserNotExist,
				Msg:  "用户不存在，请先注册",
			}, errors.New("该用户不存在")
		}
		// 不是用户不存在却还是继续出错，就说明是其他不可抗拒的因素
		return &Response{
			Code: http.StatusInternalServerError,
			Msg:  "查询数据库出现错误",
		}, nil
	}
	// 用户从数据库中找到了，检验密码
	ok, err := user.CheckPassword(req.Password)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			Msg:  "登录失败",
		}, nil
	}
	if !ok {
		return &Response{
			Code: common.PasswordErr,
			Msg:  "密码错误，登录失败",
		}, nil
	}
	// 登录成功要分发token（其他功能需要身份验证，给前端存储的）
	token, err := util.GenerateToken(uint(user.ID), user.UserName)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			Msg:  "token签发失败",
		}, nil
	}
	// 签发token后，存储到redis中
	m := map[string]string{req.UserAgent: token}
	dao.Rdb.HSet(dao.RCtx, common.AuthToken+strconv.FormatUint(uint64(user.ID), 10), m)
	return &Response{
		Code:  http.StatusOK,
		Msg:   "登录成功",
		Token: token,
	}, nil
}

func (AuthService) Register(ctx context.Context, rep *RegisterRequest) (*Response, error) {
	var user model.User
	var count int64
	dao.MysqlDB.Where(&model.User{
		UserName: rep.Username,
	}).First(&user).Count(&count)
	if count == 1 {
		return &Response{
			Code: common.UserHaveBeenRegister,
			Msg:  "用户已经注册过了",
		}, nil
	}
	// 如果数据库中没有该用户，就开始注册
	user.UserName = rep.Username
	worker, _ := util.NewWorker(0)
	id := worker.GetId()
	user.ID = id
	err := user.SetPassword(rep.Password)
	if err != nil {
		return &Response{
			Code:  http.StatusInternalServerError,
			Msg:   "数据库插入错误",
			Token: "",
		}, nil
	}
	// 加密成功就可以创建用户了
	err = dao.MysqlDB.Create(&user).Error
	if err != nil {
		return &Response{
			Code:  http.StatusInternalServerError,
			Msg:   "数据库添加数据出错",
			Token: "",
		}, nil
	}
	return &Response{
		Code: http.StatusOK,
		Msg:  "用户注册成功",
	}, nil
}

func (AuthService) Phone(ctx context.Context, req *PhoneRequest) (*Response, error) {
	phone := req.Phone
	userAgent := req.UserAgent
	code := req.Code // 验证码
	// 我们首先验证验证码，如果验证码都错误了，那么后续的工作就是免谈的
	result, err := dao.Rdb.Get(dao.RCtx, common.AuthCode+phone).Result()
	if errors.Is(err, redis.Nil) {
		// 说明没有所谓的验证码
		return &Response{
			Code: http.StatusForbidden,
			Msg:  "请输入验证码",
		}, nil
	}
	// 到这里就说明是有验证码的
	if result != code {
		// 如果验证码错误的话
		return &Response{
			Code: common.CodeErr,
			Msg:  "验证码错误",
		}, nil
	}
	var user model.User
	// 到这里说明验证码是没有任何问题的，我们就开始验证此用户是否存在
	err = dao.MysqlDB.Where(&model.User{Phone: phone}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 说明这个用户不存在，我们给它注册一个新账号
			worker, _ := util.NewWorker(0)
			id := worker.GetId()
			user.ID = id
			user.Phone = phone
			user.UserName = common.UserNamePrefix + phone // 这里就直接把用户名随机的进行定义
			err = dao.MysqlDB.Create(&user).Error
			if err != nil {
				return &Response{
					Code: http.StatusInternalServerError,
					Msg:  "数据库添加数据出错",
				}, nil
			}
		}
		return &Response{
			Code: http.StatusInternalServerError,
			Msg:  "查询数据库出现错误",
		}, nil
	}
	// 登录成功，分发token
	token, err := util.GenerateToken(uint(user.ID), user.UserName)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			Msg:  "token签发失败",
		}, nil
	}
	// 签发token之后存储到redis中
	m := map[string]string{userAgent: token}
	dao.Rdb.HSet(dao.RCtx, strconv.FormatUint(uint64(user.ID), 10), m)
	// 到这里的时候，一定是已经存在账户了，那么我们就直接登录成功就可以了
	// 注意，这里需要删除验证码！否则我们的用户登录之后马上退出登录，然后用同一个验证码是可以成功的
	dao.Rdb.Del(dao.RCtx, common.AuthCode+phone)
	return &Response{
		Code: http.StatusOK,
		Msg:  "登录成功",
	}, nil
}

func (AuthService) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	// 把redis里面的数据删除就可以了，其他的前端会给我们解决的
	userAgent := req.UserAgent
	// 我们从JWT里面来解析用户的ID
	token := req.Token
	parseToken, err := util.ParseToken(token)
	if err != nil {
		log.Println("token解析失败")
	}
	uid := parseToken.Uid
	err = dao.Rdb.HDel(dao.RCtx, common.AuthToken+strconv.Itoa(int(uid)), userAgent).Err()
	if err != nil {
		log.Println(err)
	}
	return &LogoutResponse{Code: http.StatusOK}, nil
}

func (AuthService) mustEmbedUnimplementedAuthServiceServer() {
}
