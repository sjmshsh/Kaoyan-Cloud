package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sjmshsh/grpc-gin-admin/project_api/config"
)

// RCtx 全局redis ctx
var RCtx = context.Background()
var Rdb *redis.Client

func InitRedis() {
	rdb := redis.NewClient(config.C.ReadRedisConfig())
	Rdb = rdb
}
