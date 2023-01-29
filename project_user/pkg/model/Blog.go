package model

import "time"

type Blog struct {
	Id         int64
	UserId     int64
	Content    string
	CreateTime time.Time
}

func (Blog) TableName() string {
	return "blog"
}
