package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sjmshsh/grpc-gin-admin/project_auth/config"
)

var Rdb *redis.Client
var RCtx = context.Background()

func InitRedis() {
	rdb := redis.NewClient(config.C.ReadRedisConfig())
	Rdb = rdb
}
