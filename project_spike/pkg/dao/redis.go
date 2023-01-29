package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	config2 "github.com/sjmshsh/grpc-gin-admin/project_spike/config"
)

var Rdb *redis.Client
var RCtx = context.Background()

var script = `
local redpackId = ARGV[1]
local userId = ARGV[2]
local redpackKey = KEYS[1] .. redpackId
local redpackList = KEYS[2] .. redpackId
local redpackSet = KEYS[3] .. redpackId
local res = redis.call('llen', redpackKey)
if (tonumber(res) <= 0) then
    -- 超卖了，没有库存了
    return 1
end
-- 记录用户已经买过的信息
redis.call('sadd', redpackSet, userId)
-- 给用户发一个红包并且记录下来
local money = redis.call('lpop', redpackKey)
redis.call('lpush', redpackList, tostring(money))
return 0
`
var LuaHash string

func InitRedis() {
	rdb := redis.NewClient(config2.C.ReadRedisConfig())
	Rdb = rdb
	LuaHash, _ = Rdb.ScriptLoad(context.Background(), script).Result()
}
