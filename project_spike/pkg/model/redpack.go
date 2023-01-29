package model

import "time"

type Redpacklist struct {
	Id         int64   // 作为索引
	UserId     int64   // 用户ID
	RedpackId  int64   // 红包ID
	Money      float64 // 金额
	CreateTime time.Time
}

type Redpack struct {
	Id         int64
	RedpackId  int64
	Amount     float64
	count      int64
	CreateTime time.Time
	UserId     int64
}
