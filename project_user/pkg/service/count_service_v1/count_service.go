package count_service_v1

import (
	"context"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/util"
	"github.com/sjmshsh/grpc-gin-admin/project_count_service/pkg/common/redisConstance"
	"github.com/sjmshsh/grpc-gin-admin/project_count_service/pkg/dao"
	"github.com/sjmshsh/grpc-gin-admin/project_count_service/pkg/model"
	"log"
	"strconv"
	"time"
)

type CountService struct {
	UnimplementedCountServiceServer
}

func New() *CountService {
	return &CountService{}
}

func (c *CountService) Count(ctx context.Context, req *CountRequest) (*Response, error) {
	// 文章ID
	id := req.Id
	// 用户ID
	// 类型,是点赞还是评论等等
	ty := req.Type
	symbol := req.Symbol
	if ty == redisConstance.StarType || ty == redisConstance.WatchType {
		// 如果类型是点赞数或者浏览数的话，就说明是模糊计数
		fuzzyCount(id, int(ty), int(symbol))
	}
	// 这里是可以一直拓展的
	if ty == redisConstance.CommentType || ty == redisConstance.RepostType {
		// 类型是评论类型或者转发类型，就说明是精准计数
		preciseCount(id, int(ty), int(symbol))
	}
	return &Response{
		Data: id,
	}, nil
}

// 我事先规定好，模糊计数有：点赞，浏览，就这两个，其他的全部是精准计数
// 当 key 不存在时，返回 -2 。
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
// 否则，以秒为单位，返回 key 的剩余生存时间。
// 模糊计数
// 读Redis，写Redis，Cron定期同步检测Redis和DB，但是频率很低
// 每过一段时间就使用定时器把Redis里面的数据写入到DB中
// 返回的是计数之后的数量值，直接返回给前端
func fuzzyCount(id int64, ty, symbol int) {
	i := strconv.FormatInt(id, 10)
	key := redisConstance.CounterFuzzyService + i
	// 先查看这个key是否存在，如果不存在的话需要创建
	duration, err := dao.Rdb.TTL(dao.RCtx, key).Result()
	if err != nil {
		log.Println(err)
	}
	if duration == -2 {
		var fuzzycount []model.Fuzzy
		// 说明缓存直接不存在，我们需要创建缓存，构建缓存的时候就需要查询数据库, 前提是我们的数据库已经构建了，所以这里我们还需要去判断一下
		err := dao.MysqlDB.Debug().Raw("select * from fuzzy where cid = ?", id).Scan(&fuzzycount).Error
		if err != nil {
			log.Println(err)
		}
		if fuzzycount == nil {
			// 说明数据库里面还没有，我们要在数据库里面新创建
			worker, err := util.NewWorker(10)
			if err != nil {
				log.Println(err)
			}
			// 如果是数据库里面查不到的话，那么我们就需要插入数据库，首次插入，要么是0，要么是1
			err = dao.MysqlDB.Debug().
				Exec("insert into fuzzy (id, `type`, num, cid) values(?, ?, ?, ?)",
					worker.GetId(), ty, 1, id).Error
			if err != nil {
				log.Println(err)
			}
			if ty == redisConstance.StarType {
				err = dao.MysqlDB.Debug().
					Exec("insert into fuzzy (id, `type`, num, cid) values(?, ?, ?, ?)",
						worker.GetId(), redisConstance.WatchType, 0, id).Error
			}
			if ty == redisConstance.WatchType {
				err = dao.MysqlDB.Debug().
					Exec("insert into fuzzy (id, `type`, num, cid) values(?, ?, ?, ?)",
						worker.GetId(), redisConstance.StarType, 0, id).Error
			}
			if err != nil {
				log.Println(err)
			}
			// 构建缓存
			if ty == redisConstance.StarType {
				dao.Rdb.HSet(dao.RCtx, key, redisConstance.StarNum, 0)
				dao.Rdb.HSet(dao.RCtx, key, redisConstance.WatchNum, 0)
			} else {
				dao.Rdb.HSet(dao.RCtx, key, redisConstance.WatchNum, 0)
				dao.Rdb.HSet(dao.RCtx, key, redisConstance.StarNum, 0)
			}
			// 一天就过期，对大V特殊处理
			dao.Rdb.Expire(dao.RCtx, key, time.Minute*60*24)
		} else {
			// 说明数据库里面已经有数据了，我们根据数据库里面的数据来构建缓存
			for i := 0; i < len(fuzzycount); i++ {
				var value string
				if fuzzycount[i].Type == 1 {
					value = redisConstance.StarNum
				}
				if fuzzycount[i].Type == 2 {
					value = redisConstance.WatchNum
				}
				dao.Rdb.HSet(dao.RCtx, key, value, fuzzycount[i].Num)
			}
			dao.Rdb.Expire(dao.RCtx, key, time.Minute*60*24)
		}
		duration = time.Minute * 60 * 24
	}
	if duration < redisConstance.CounterServiceTTL/3 {
		// 如果缓存的生命周期已经低于1/3，就更新缓存过期时间
		dao.Rdb.Expire(dao.RCtx, key, time.Minute*60*24)
	}
	var value string

	// 1代表点赞
	if ty == redisConstance.StarType {
		value = redisConstance.StarNum
	}
	if ty == redisConstance.WatchType {
		value = redisConstance.WatchNum
	}
	// 这里可以继续加浏览量等等，这里只写了一个而已
	result, err := dao.Rdb.HGet(dao.RCtx, key, value).Result()
	if err != nil {
		log.Println(err)
	}
	r, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		log.Println(err)
	}
	var count int
	if symbol == 1 {
		// 1表示+
		count = int(r + 1)
	} else {
		count = int(r - 1)
		if count < 0 {
			count = 0
		}
	}
	// 然后把这个计数++
	dao.Rdb.HSet(dao.RCtx, key, value, count)
}

// 写的时候删除缓存，查的时候构建缓存
// 使用先更新数据库，再删除缓存的策略
// 读Redis，写的时候写入DB，更新Redis，定期同步检测Redis和DB，高频率，做到一定精准，同时消息队列
// 异步写入DB的时候要做好一切措施保证消息的最终一致性
func preciseCount(id int64, ty, symbol int) {
	// 我先得到我目前的评论数量
	i := strconv.FormatInt(id, 10)
	key := redisConstance.CountPreciseService + i
	// 我们首先要确认数据库里面有没有数据
	var res []model.Precise
	// TODO 这里有点儿耗时了
	err := dao.MysqlDB.Debug().Raw("select * from precise where cid = ?", id).Scan(&res).Error
	if err != nil {
		log.Println(err)
	}
	if res == nil {
		// 如果结果是空，说明数据库里面根本没有数据
		// 我们此时就需要把精准查询的数据创建出来
		worker, err := util.NewWorker(11)
		if err != nil {
			log.Println(err)
		}
		if ty == redisConstance.CommentType {
			err = dao.MysqlDB.Debug().
				Exec("insert into precise (id, `type`, num, cid) values (?, ?, ?, ?)",
					worker.GetId(), ty, 1, id).Error
			if err != nil {
				log.Println(err)
			}
			err = dao.MysqlDB.Debug().
				Exec("insert into precise (id, `type`, num, cid) values (?, ?, ?, ?)",
					worker.GetId(), redisConstance.RepostType, 0, id).Error
			if err != nil {
				log.Println(err)
			}
		} else if ty == redisConstance.RepostType {
			err = dao.MysqlDB.Debug().
				Exec("insert into precise (id, `type`, num, cid) values (?, ?, ?, ?)",
					worker.GetId(), ty, 1, id).Error
			if err != nil {
				log.Println(err)
			}
			err = dao.MysqlDB.Debug().
				Exec("insert into precise (id, `type`, num, cid) values (?, ?, ?, ?)",
					worker.GetId(), redisConstance.CommentType, 0, id).Error
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		// 查询结果不为空
		for i := 0; i < len(res); i++ {
			if res[i].Type == ty {
				if symbol == 1 {
					dao.MysqlDB.Debug().Exec("update precise set num = ? where id = ?", res[i].Num+1, res[i].Id)
				} else {
					if res[i].Num-1 >= 0 {
						dao.MysqlDB.Debug().Exec("update precise set num = ? where id = ?", res[i].Num-1, res[i].Id)
					}
				}
			}
		}
	}
	// 然后把这个缓存直接删除, 如果这个没有构建缓存的话更好, 这里进行异步删除，把这个东西放入到消息队列里面去
	// 如果删除失败了，那么那么就在消息队列中进行重试，这里就要又涉及到了消息队列的相关知识了
	err = dao.Rdb.Del(dao.RCtx, key).Err()
	if err != nil {
		// 如果程序走到这里了就说明删除缓存失败了，我这个时候应该放到消息队列里面让它持续的进行删除
		mq := dao.NewRabbitMQTopics(redisConstance.RetryMessageDeleteFailedQueueName,
			redisConstance.RetryMessageDeleteFailedExchangeName,
			redisConstance.RetryMessageDeleteFailedRoute)
		mq.PublishTopics(key)
	}
}

// ConsumeDeleteFailed 这里不需要处理消息的幂等性，因为我多删除了一次也无所谓，这个消息不具有幂等属性
func ConsumeDeleteFailed() {
	mq := dao.NewRabbitMQTopics(redisConstance.RetryMessageDeleteFailedQueueName,
		redisConstance.RetryMessageDeleteFailedExchangeName,
		redisConstance.RetryMessageDeleteFailedRoute)
	go mq.ConsumeFailedCache()
}

// InitFuzzyScheduler 这边每过24小时去校验一次
// InitFuzzyScheduler 点赞数量可以属于模糊查询，所以我们在这个定时器里面定期定时这个东西
// InitFuzzyScheduler 初始化模糊查询定时器
func InitFuzzyScheduler() {
	ticker := time.NewTicker(time.Minute)
	go func() { // 用新协程去执行定时认为
		// 思考，这里是不是应该睡一会儿，这样一直循环对CPU损伤太大了
		time.Sleep(time.Second * 17) // 这里睡眠的时间有讲究，不能太长也不能太短
		for {
			// 用一个死循环不定的执行，否则只会执行一次
			select {
			case <-ticker.C: // 时间到了就会执行这个分支的代码
				FuzzyScheduler()
			}
		}
	}()
}

// TODO 写一个定时器，定时的寻找大V用户，如果关注着数量大于10000的话，那么我们就认为是一个大V

func FuzzyScheduler() {
	// 查询数据库，得到数据库里面所有type = 1的东西，用于批量更新Redis
	// TODO 当然，这是一个巨大的查询，因此我们需要分批来做，避免MySQL突然宕机了
	var res []model.Fuzzy
	err := dao.MysqlDB.Debug().
		Raw("select * from fuzzy").
		Scan(&res).Error
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(res); i++ {
		id := strconv.FormatInt(res[i].Cid, 10)
		key := redisConstance.CounterFuzzyService + id
		if res[i].Type == redisConstance.StarType {
			result, err := dao.Rdb.HGet(dao.RCtx, key, redisConstance.StarNum).Result()
			// 这里有一个巨大的问题，就是如果缓存这个时候已经过期了，那么此时我HGet之后返回的是0，然后我再更新数据库，数据库就被清零了！！
			if err != nil {
				log.Println(err)
			}
			if result != "" {
				dao.MysqlDB.Debug().Exec("update fuzzy set num = ? where id = ?", result, res[i].Id)
			}
		} else if res[i].Type == redisConstance.WatchType {
			result, err := dao.Rdb.HGet(dao.RCtx, key, redisConstance.WatchNum).Result()
			if err != nil {
				log.Println(err)
			}
			if result != "" {
				dao.MysqlDB.Debug().Exec("update fuzzy set num = ? where id = ?", result, res[i].Id)
			}
		}
	}
}

func (c *CountService) GetCount(ctx context.Context, req *CountRequest) (*Response, error) {
	i := req.Id
	id := strconv.FormatInt(i, 10)
	// 类型,是点赞还是评论等等
	ty := int(req.Type)
	// 先查Redis
	var key string
	var num int
	if ty == redisConstance.StarType || ty == redisConstance.WatchType {
		// 模糊计数
		key = redisConstance.CounterFuzzyService + id
		var value string
		if ty == redisConstance.StarType {
			value = redisConstance.StarNum
		} else if ty == redisConstance.WatchType {
			value = redisConstance.WatchNum
		}
		result, err := dao.Rdb.HGet(dao.RCtx, key, value).Result()
		if err != nil {
			log.Println(err)
		}
		if result != "" {
			parseInt, err := strconv.ParseInt(result, 10, 64)
			if err != nil {
				log.Println(err)
			}
			num = int(parseInt)
		} else {
			var res []model.Fuzzy
			// 说明缓存里面没有数据，要从数据库里面去读, 读取完毕了之后我们顺便构建缓存
			err := dao.MysqlDB.Debug().Raw("select * from fuzzy where cid = ?", id).Scan(&res).Error
			if err != nil {
				log.Println(err)
			}
			// 构建缓存
			for i := 0; i < len(res); i++ {
				if res[i].Type == redisConstance.StarType {
					dao.Rdb.HSet(dao.RCtx, key, redisConstance.StarNum, res[i].Num)
				} else if res[i].Type == redisConstance.WatchType {
					dao.Rdb.HSet(dao.RCtx, key, redisConstance.WatchNum, res[i].Num)
				}
				if res[i].Type == ty {
					num = res[i].Num
				}
			}
			// 设置一天的过期时间
			dao.Rdb.Expire(dao.RCtx, key, time.Minute*60*24)
		}
	} else if ty == redisConstance.RepostType || ty == redisConstance.CommentType {
		// 精准计数
		key = redisConstance.CountPreciseService + id
		var value string
		if ty == redisConstance.RepostType {
			value = redisConstance.RepostNum
		} else if ty == redisConstance.CommentType {
			value = redisConstance.CommentNum
		}
		result, err := dao.Rdb.HGet(dao.RCtx, key, value).Result()
		if err != nil {
			log.Println(err)
		}
		if result != "" {
			parseInt, err := strconv.ParseInt(result, 10, 64)
			if err != nil {
				log.Println(err)
			}
			num = int(parseInt)
		} else {
			var res []model.Precise
			// 说明缓存里面没有数据，要从数据库里面去读, 读取完毕了之后我们顺便构建缓存
			err := dao.MysqlDB.Debug().Raw("select * from precise where cid = ?", id).Scan(&res).Error
			if err != nil {
				log.Println(err)
			}
			// 构建缓存
			for i := 0; i < len(res); i++ {
				if res[i].Type == redisConstance.RepostType {
					dao.Rdb.HSet(dao.RCtx, key, redisConstance.RepostNum, res[i].Num)
				} else if res[i].Type == redisConstance.CommentType {
					dao.Rdb.HSet(dao.RCtx, key, redisConstance.CommentNum, res[i].Num)
				}
				if res[i].Type == ty {
					num = res[i].Num
				}
			}
			// 设置一天的过期时间
			dao.Rdb.Expire(dao.RCtx, key, time.Minute*60*24)
		}
	}
	return &Response{
		Data: int64(num),
	}, nil
}
