package user_service_v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sjmshsh/grpc-gin-admin/project_user/pkg/common"
	"github.com/sjmshsh/grpc-gin-admin/project_user/pkg/dao"
	"github.com/sjmshsh/grpc-gin-admin/project_user/pkg/model"
	"github.com/sjmshsh/grpc-gin-admin/project_user/pkg/service/count_service_v1"
	"github.com/sjmshsh/grpc-gin-admin/project_user/pkg/util"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type UserService struct {
	UnimplementedUserServiceServer
}

func New() *UserService {
	return &UserService{}
}

type VUser struct {
	Follower  int64  // 粉丝数量
	Attention int64  // 关注数量
	Name      string // 大V的用户名
}

// LocalCache 定义一个二级缓存
//var LocalCache = make(map[string]User)

func (u *UserService) CheckFileMd5(ctx context.Context, req *CheckFileMd5Request) (*CheckFileMd5Response, error) {
	md5 := req.Md5
	chunks := req.Chunks
	// 从redis里面获取这个文件的相关信息
	status, err := dao.Rdb.HGet(dao.RCtx, common.FileUploadStatus, md5).Result()
	if errors.Is(err, redis.Nil) {
		// 如果没有找到的话
		// 说明文件还没有开始上传
		return &CheckFileMd5Response{
			Flg:           common.NoHave,
			MisschunkList: nil,
		}, nil
	}
	flg, _ := strconv.ParseBool(status)
	if flg == true {
		// 这样说明文件已经上传完成了
		return &CheckFileMd5Response{
			Flg:           common.IsHave,
			MisschunkList: nil,
		}, nil
	}
	// 到这里说明文件还没有上传完成，我们使用位图的方式判断哪些东西没有上传完成
	key := common.FileProcessingStatus + md5
	missChunkList := make([]int32, 0)
	for i := 0; i < int(chunks); i++ {
		// 偏移量是多少就代表是多少号分片
		result, err := dao.Rdb.GetBit(dao.RCtx, key, int64(i)).Result()
		if err != nil {
			log.Println(err)
		}
		if result == 0 {
			// 如果这个编号是0，代表还没有上传成功
			missChunkList = append(missChunkList, int32(i))
		}
		// 如果编号是1的话，就说明是上传成功的
	}
	return &CheckFileMd5Response{
		Flg:           common.IngHave,
		MisschunkList: missChunkList,
	}, nil
}

func (u *UserService) UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	// 这里做的工作主要就是入库了
	file := model.File{
		Id:       req.Id,
		Md5:      req.Md5,
		Name:     req.Name,
		Size:     req.Size,
		Addr:     req.Path,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		Status:   1,
	}
	err := dao.MysqlDB.Create(&file).Error
	if err != nil {
		log.Println(err)
	}
	return &UploadFileResponse{
		Code: http.StatusOK,
		Msg:  "文件上传成功",
	}, nil
}

func (u *UserService) GetSign(ctx context.Context, request *GetSignRequest) (*Response, error) {
	userId := request.UserId
	id := strconv.FormatInt(userId, 10)
	year := request.Year
	month := request.Month
	// 拼接redis的key
	key := common.UserCheckIn + id + ":" + year + ":" + month
	fmt.Println(key)
	// 通过bitfield命令返回整个的数组
	// 数组的第一个元素就是一个int64类型的值，我们通过位运算进行操作
	s := fmt.Sprintf("i%d", 31)
	fmt.Println(s)
	result, err := dao.Rdb.BitField(dao.RCtx, key, "get", s, 0).Result()
	if err != nil {
		log.Println(err)
	}
	num := result[0]
	fmt.Println(num)
	arr := make([]int64, 31)
	for i := 0; i < 31; i++ {
		// 让这个数字与1做与运算，得到数据的最后一个比特
		if (num & 1) == 0 {
			// 如果为0，说明未签到
			arr[i] = 0

		} else {
			// 如果不为0，说明已经签到了，计数器+1
			arr[i] = 1
		}
		// 把数字右移动一位，抛弃最后一个bit位，继续下一个bit位
		num = num >> 1
	}
	return &Response{
		Status: http.StatusOK,
		Msg:    "获取签到信息成功",
		Data:   arr,
	}, nil
}

func (u *UserService) WatchUv(ctx context.Context, request *WatchUvRequest) (*WatchUvResponse, error) {
	result, err := dao.Rdb.PFCount(dao.RCtx, common.WebUv).Result()
	if err != nil {
		log.Println(err)
	}
	return &WatchUvResponse{
		Uv: result,
	}, nil
}

func (u *UserService) CheckIn(ctx context.Context, request *CheckSignRequest) (*Response, error) {
	userId := request.UserId
	year := request.Year
	month := request.Month
	day := request.Day
	d, _ := strconv.ParseInt(day, 10, 64)
	// 组装redis的key
	id := strconv.FormatInt(userId, 10)
	key := common.UserCheckIn + id
	// 拼装
	// 2023:01:15  2023 01 15
	value := fmt.Sprintf(":%s:%s", year, month)
	key = key + value
	// 签到的代码
	dao.Rdb.SetBit(dao.RCtx, key, d+1, 1)
	// 设置过期时间, 30天，可以长一点
	dao.Rdb.Expire(dao.RCtx, key, time.Hour*24*30)
	// 缓存层已经设置，接下来使用消息队列异步存储到存储层MySQL
	message := model.CheckInMessage{
		UserId: userId,
		Year:   year,
		Month:  month,
		Day:    day,
	}
	mq := dao.NewRabbitMQTopics("sign", "sign-", "hello")
	mq.PublishTopics(message)
	return &Response{
		Status: http.StatusOK,
		Msg:    "用户签到成功",
	}, nil
}

func (u *UserService) Location(ctx context.Context, request *LocationRequest) (*Response, error) {
	longitude := request.Longitude
	o, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		log.Println(err)
	}
	latitude := request.Latitude
	a, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		log.Println(err)
	}
	userId := request.UserId
	location := request.Location
	s := strconv.FormatInt(userId, 10)
	key := common.UserLocation + s
	dao.Rdb.GeoAdd(dao.RCtx, key, &redis.GeoLocation{
		Name:      location,
		Longitude: o,
		Latitude:  a,
	})
	return &Response{
		Status: http.StatusOK,
		Msg:    "插入信息成功",
	}, nil
}

func (u *UserService) FindFriend(ctx context.Context, request *FindFriendRequest) (*FindFriendResponse, error) {
	longitude := request.Longitude
	latitude := request.Latitude
	o, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		log.Println(err)
	}
	a, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		log.Println(err)
	}
	key := common.UserLocation
	result, err := dao.Rdb.GeoRadius(dao.RCtx, key, o, a, &redis.GeoRadiusQuery{
		Radius:      5,
		Unit:        "km",
		WithCoord:   false, //传入WITHCOORD参数，则返回结果会带上匹配位置的经纬度
		WithDist:    true,  //传入WITHDIST参数，则返回结果会带上匹配位置与给定地理位置的距离。
		WithGeoHash: false, //传入WITHHASH参数，则返回结果会带上匹配位置的hash值
		Sort:        "ASC", //默认结果是未排序的，传入ASC为从近到远排序，传入DESC为从远到近排序。
	}).Result()
	if err != nil {
		log.Println(err)
	}
	N := len(result)
	name := make([]string, N)
	dist := make([]float32, N)
	for i := 0; i < N; i++ {
		name[i] = result[i].Name
		dist[i] = float32(result[i].Dist)
	}
	return &FindFriendResponse{
		Name: name,
		Dist: dist,
	}, nil
}

// PostBlog 发博客的时候需要考虑Feed流的问题
func (u *UserService) PostBlog(ctx context.Context, req *PostBlogRequest) (*Response, error) {
	content := req.Content
	userId := req.UserId
	// 往数据库里面插入博客
	worker, err := util.NewWorker(1)
	if err != nil {
		log.Println(err)
	}
	id := worker.GetId()
	blog := &model.Blog{
		Id:         id,
		UserId:     userId,
		Content:    content,
		CreateTime: time.Now(),
	}
	dao.MysqlDB.Debug().Create(blog)
	// 然后我们把博客发到feed里面去
	// 首先得到我所有的跟随着
	var follower []int64
	dao.MysqlDB.Raw("select follower_id from follower where user_id = ?", userId).Scan(&follower)
	// 得到了所有我追随者的id
	for i := 0; i < len(follower); i++ {
		sid := strconv.FormatInt(follower[i], 10)
		// 为追随者构建feed流
		key := common.Feed + sid
		dao.Rdb.LPush(dao.RCtx, key, id)
	}
	formatInt := strconv.FormatInt(id, 10)
	return &Response{
		Status: http.StatusOK,
		Msg:    formatInt,
	}, nil
}

// Watch 关注，不走计数服务
// 走计数服务的话会导致整个的代码紊乱，很难处理
// 普通关注了之后把缓存删了
// 大V用户关注了之后对缓存进行更改，异步更新数据库
func (u *UserService) Watch(ctx context.Context, req *WatchRequest) (*Response, error) {
	userId := req.UserId
	userid := strconv.FormatInt(userId, 10)
	attentionUserId := req.AttentionUserId
	a := strconv.FormatInt(attentionUserId, 10)
	worker, err := util.NewWorker(1)
	if err != nil {
		log.Println(err)
	}
	aid := worker.GetId()
	// 先插入数据库
	attention := &model.Attention{
		Id:          aid,
		UserId:      userId,
		AttentionId: attentionUserId,
		Flg:         1, // 代表已经关注了
	}
	fid := worker.GetId()
	follower := &model.Follower{
		Id:         fid,
		UserId:     attentionUserId,
		FollowerId: userId,
	}
	// 首先把用户关注数量的缓存给删了，这是精准缓存
	group := sync.WaitGroup{}
	group.Add(4)
	go func() {
		dao.MysqlDB.Debug().Create(attention)
		group.Done()
	}()
	go func() {
		dao.MysqlDB.Debug().Create(follower)
		group.Done()
	}()
	go func() {
		key := common.EXIST + userid
		dao.Rdb.HSet(dao.RCtx, key, common.ATTENTION, 1)
		group.Done()
	}()
	go func() {
		key1 := common.AttentionList + userid
		key2 := common.FollowerList + a
		// 把缓存删除掉
		dao.Rdb.Del(dao.RCtx, key1)
		dao.Rdb.Del(dao.RCtx, key2)
		group.Done()
	}()
	group.Wait()
	return &Response{
		Msg: "已关注了",
	}, nil
}

func (u *UserService) Comment(ctx context.Context, req *CommentRequest) (*CommentResponse, error) {
	content := req.Content
	id := req.Id
	// 插入到数据库里
	worder, err := util.NewWorker(1)
	if err != nil {
		log.Println(err)
	}
	getId := worder.GetId()
	comment := &model.Comment{
		Id:       getId,
		Type:     common.BlogComment,
		MemberId: id,
		State:    1, // 目前还没有进行敏感词匹配
		Content:  content,
	}
	dao.MysqlDB.Debug().Create(comment)
	return &CommentResponse{
		Msg: "评论成功",
	}, nil
}

type R struct {
	Content  string // 博客内容
	Likes    int    // 点赞数量
	Comments int    // 评论数量
}

// GetFeedList 有博客内容，用户相关信息
func (u *UserService) GetFeedList(ctx context.Context, req *GetFeedListRequest) (*GetFeedListResponse, error) {
	// 获取缓存里面的信息，如果缓存里面找不到就说明目前没有推荐，就不推荐
	userId := req.UserId
	id := strconv.FormatInt(userId, 10)
	// 这个里面就是json串
	result, err := dao.Rdb.LRange(dao.RCtx, common.Feed+id, req.Start, req.Start+req.Offset).Result()
	if err != nil {
		log.Println(err)
	}
	var feedlist []string
	// 这个result里面存放的是我的feed的id，也就是文章的id
	for i := 0; i < len(result); i++ {
		// 这个时候我们需要先获取博客的内容，然后获取博客的点赞数量，博客的评论数量
		// 博客内容
		// 先从缓存里面获取博客内容
		key := common.Content + result[i]
		r, err := dao.Rdb.Get(dao.RCtx, key).Result()
		if err != nil {
			log.Println(err)
		}
		if r == "" {
			// 是空，说明缓存里面没有，我们需要从MySQL里面获取缓存内容
			err := dao.MysqlDB.Debug().Raw("select content from blog where id = ?", result[i]).Scan(&r).Error
			if err != nil {
				log.Println(err)
			}
			dao.Rdb.Set(dao.RCtx, common.Content+result[i], r, time.Minute*60*24)
		}
		parseInt, err := strconv.ParseInt(result[i], 10, 64)
		if err != nil {
			log.Println(err)
		}
		group := sync.WaitGroup{}
		group.Add(2)
		var resp1 *count_service_v1.Response
		var resp2 *count_service_v1.Response
		go func() {
			// 点赞数量和评论数量，调用计数服务
			resp1, err = CountServieClient.GetCount(context.Background(), &count_service_v1.CountRequest{
				Type: common.StarType,
				Id:   parseInt,
			})
			if err != nil {
				log.Println(err)
			}
			group.Done()
		}()
		go func() {
			resp2, err = CountServieClient.GetCount(context.Background(), &count_service_v1.CountRequest{
				Type: common.CommentType,
				Id:   parseInt,
			})
			if err != nil {
				log.Println(err)
			}
			group.Done()
		}()
		group.Wait()
		res := &R{
			Content:  r,
			Likes:    int(resp1.Data),
			Comments: int(resp2.Data),
		}
		marshal, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
		}
		feedlist = append(feedlist, string(marshal))
	}
	return &GetFeedListResponse{
		Feed: feedlist,
	}, nil
}

func (u *UserService) OList(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	t := req.Type
	id := req.UserId
	offset := req.Offset
	start := req.Start
	userId := strconv.FormatInt(id, 10)
	var response []string
	if t == common.Like {
		// 点赞列表
		// TODO 其实代码都是差不多的，这里暂时省略，后期有时间再写
	} else if t == common.Follower {
		// 说明是关注者列表
		// 从缓存中获取
	FLAG:
		key := common.FollowerList + userId
		result, err := dao.Rdb.LRange(dao.RCtx, key, start, start+offset).Result()
		if err != nil {
			log.Println(err)
		}
		// result 里面存放的是用户的ID
		// 我们根据用户ID可以得到用户的相关信息
		if len(result) == 0 {
			// 构建缓存的时候使用互斥锁
			formatInt := strconv.FormatInt(id, 10)
			lockKey := common.LockListKey + formatInt
			// 如果获取锁失败的话
			lock := util.TryLock(lockKey)
			if !lock {
				// 没有获取到锁, 直接睡0.5秒
				time.Sleep(time.Millisecond * 500)
				goto FLAG
			}
			// 说明缓存里面是没有东西的，我们需要去数据库里面查询
			err := dao.MysqlDB.Debug().Raw("select follower_id from attention where user_id = ?", userId).Scan(&result).Error
			if err != nil {
				log.Println(err)
			}
			// 然后构建缓存
			dao.Rdb.LPush(dao.RCtx, key, result)
			dao.Rdb.Expire(dao.RCtx, key, time.Minute*60*24)
		}
		// 到这里为止result里面已经百分百有东西了
		for i := 0; i < len(result); i++ {
			// 提取出用户ID
			id := result[i]
			// 目前只有姓名，后期再拓展
			key := common.User + id
			s, err := dao.Rdb.Get(dao.RCtx, key).Result()
			if err != nil {
				log.Println(err)
			}
			var l string
			if s == "" {
				// 缓存里面没有，从数据库里面查找
				err := dao.MysqlDB.Debug().Raw("select user_name from user where id = ?", id).Scan(&l).Error
				if err != nil {
					log.Println(err)
				}
				// 构建用户缓存
				dao.Rdb.Set(dao.RCtx, key, l, time.Minute*60*24+time.Duration(rand.Int()))
			}
			response = append(response, l)
		}
	}
	return &ListResponse{
		UserList: response,
	}, nil
}

func (u *UserService) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	t := req.Type
	id := req.UserId
	offset := req.Offset
	start := req.Start
	userId := strconv.FormatInt(id, 10)
	var response []string
	// 从缓存重获取粉丝列表, 缓存里面存的是粉丝ID
	if t == common.Attention {
		// 说明是关注者列表
		// 从缓存中获取
	FLAG:
		key := common.AttentionList + userId
		result, err := dao.Rdb.LRange(dao.RCtx, key, start, start+offset).Result()
		if err != nil {
			log.Println(err)
		}
		// result 里面存放的是用户的ID
		// 我们根据用户ID可以得到用户的相关信息
		if len(result) == 0 {
			// 说明缓存里面是没有东西的，我们需要去数据库里面查询
			// 构建缓存的时候使用互斥锁
			formatInt := strconv.FormatInt(id, 10)
			lockKey := common.LockListKey + formatInt
			// 如果获取锁失败的话
			lock := util.TryLock(lockKey)
			if !lock {
				// 没有获取到锁, 直接睡0.5秒
				time.Sleep(time.Millisecond * 500)
				goto FLAG
			}
			// 获取到了锁
			err := dao.MysqlDB.Debug().Raw("select attention_id from attention where user_id = ?", userId).Scan(&result).Error
			if err != nil {
				log.Println(err)
			}
			// 然后构建缓存
			dao.Rdb.LPush(dao.RCtx, key, result)
			ran := rand.Uint64()
			dao.Rdb.Expire(dao.RCtx, key, time.Minute*60*24+time.Second*time.Duration(ran))
		}
		// 到这里为止result里面已经百分百有东西了
		for i := 0; i < len(result); i++ {
			// 提取出用户ID
			id := result[i]
			// 目前只有姓名，后期再拓展
			key := common.User + id
			s, err := dao.Rdb.Get(dao.RCtx, key).Result()
			if err != nil {
				log.Println(err)
			}
			var l string
			if s == "" {
				// 缓存里面没有，从数据库里面查找
				err := dao.MysqlDB.Debug().Raw("select user_name from user where id = ?", id).Scan(&l).Error
				if err != nil {
					log.Println(err)
				}
				// 构建用户缓存
				dao.Rdb.Set(dao.RCtx, key, l, time.Minute*60*24)
			}
			response = append(response, l)
		}
	}
	return &ListResponse{
		UserList: response,
	}, nil
}

func InitSignConsumer() {
	mq := dao.NewRabbitMQTopics("sign", "sign-", "hello")
	go mq.ConsumeTopicsCheckIn()
}
