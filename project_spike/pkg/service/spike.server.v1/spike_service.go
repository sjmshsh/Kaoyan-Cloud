package spike_service_v1

import (
	"context"
	"github.com/sjmshsh/grpc-gin-admin/project_spike/pkg/code"
	"github.com/sjmshsh/grpc-gin-admin/project_spike/pkg/common"
	"github.com/sjmshsh/grpc-gin-admin/project_spike/pkg/dao"
	"github.com/sjmshsh/grpc-gin-admin/project_spike/pkg/model"
	"github.com/sjmshsh/grpc-gin-admin/project_spike/pkg/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

type SpikeService struct {
	UnimplementedSpikeServiceServer
}

func New() *SpikeService {
	return &SpikeService{}
}

// SendRedPack 发红包
func (sp *SpikeService) SendRedPack(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	// 拆解红包
	packet := util.SplitRedPacket(float64(req.Amount), req.Number)
	log.Println(packet)
	// 为红包生成全局唯一的ID
	worker, err := util.NewWorker(1)
	if err != nil {
		log.Println(err)
	}
	id := worker.GetId()
	key := common.RedPacketKey + strconv.FormatInt(id, 10)
	log.Println(key)
	// 采用list存储红包
	for i := 0; i < int(req.Number); i++ {
		dao.Rdb.LPush(context.Background(), key, packet[i])
	}
	// 设置过期时间,24小时即可
	dao.Rdb.Expire(context.Background(), key, time.Minute*60*24)
	// 把红包的相关信息添加进入数据库
	redpack := &model.Redpack{
		Id:         worker.GetId(),
		RedpackId:  id,
		Amount:     float64(req.Amount),
		CreateTime: time.Now(),
		UserId:     req.UserID,
	}
	err = dao.MysqlDB.Debug().Create(redpack).Error
	if err != nil {
		log.Println(err)
	}
	return &SendResponse{
		Code:        "400",
		Msg:         "发送成功",
		RedPacketId: id,
	}, err
}

// RecvRedPack 抢红包
func (sp *SpikeService) RecvRedPack(ctx context.Context, req *RecvRequest) (*RecvResponse, error) {
	// 获取红包ID
	id := req.Id
	// 获取用户ID
	userId := req.UserId
	key := common.RedPacketKey + strconv.FormatInt(id, 10)
	// 先查询红包是否还有库存(key是否为空)，如果没有库存了的话就直接给出抢红包的列表
	len, err := dao.Rdb.LLen(context.Background(), key).Result()
	if err != nil {
		log.Println(err)
	}
	if len == 0 {
		// 说明已经没有库存了
		return &RecvResponse{
			Code: code.InventoryShortage,
			Msg:  "库存不足",
		}, nil
	}
	// 到这里说明还有库存, 用户可用点开进行抢红包
	// 先判断是否还有库存，如果有库存的话再判断
	result, err := dao.Rdb.EvalSha(context.Background(), dao.LuaHash, []string{common.RedPacketKey, common.RedPackList, common.RedPackSet}, id, userId).Result()
	if err != nil {
		log.Println(err)
	}
	res := result.(int64)
	log.Println("----------------")
	log.Println(res)
	if res == 1 {
		return &RecvResponse{
			Code: code.StackShortage,
			Msg:  "没有库存了",
		}, nil
	}
	if res == 2 {
		return &RecvResponse{
			Code: code.UserHaveBought,
			Msg:  "用户已经买过了",
		}, nil
	}
	// 返回0说明成功了
	// 不用入库，24小时之后直接就无法查看了
	return &RecvResponse{
		Code: http.StatusOK,
		Msg:  "抢红包成功",
	}, nil
}

func (sp *SpikeService) ListRedPack(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	// 先查看缓存里面有没有
	// 有的话就从缓存里面取
	// 如果缓存里面没有说明已经过去了24小时了，就直接寄

	// 获取红包ID
	id := req.Id
	s := strconv.FormatInt(id, 10)
	// 我是从左侧压入的，因此我先抢到红包的用户是在最左边，所以我们弹出的时候也从最左边弹出
	// 先获取列表的长度
	key := common.RedPackList + s
	len, err := dao.Rdb.LLen(dao.RCtx, key).Result()
	if err != nil {
		log.Println(err)
	}
	list, err := dao.Rdb.LRange(dao.RCtx, key, 0, len-1).Result()
	if err != nil {
		log.Println(err)
	}
	return &ListResponse{
		List: list,
	}, nil
}

func (SpikeService) mustEmbedUnimplementedSpikeServiceServer() {
}
