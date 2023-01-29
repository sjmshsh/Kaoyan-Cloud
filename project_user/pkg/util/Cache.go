package util

import (
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/dao"
	"time"
)

func TryLock(key string) bool {
	// 设置10秒的过期时间
	err := dao.Rdb.Set(dao.RCtx, key, "1", time.Second*10).Err()
	if err != nil {
		return false
	}
	return true
}

func UnLock(key string) {
	dao.Rdb.Del(dao.RCtx, key)
}
