package model

import "time"

type File struct {
	Id       int64
	Md5      string
	Name     string
	Size     int64
	Addr     string
	CreateAt time.Time
	UpdateAt time.Time
	// 0文件已经被删除了，1表示文件还没有被删除
	Status int
}

func (File) TableName() string {
	return "file"
}

// Sign 签到
type Sign struct {
	Id     int64  `json:"id"`
	UserId int64  `json:"user_id"`
	Year   string `json:"year"`
	Month  string `json:"month"`
	Day    string `json:"day"`
}

func (Sign) TableName() string {
	return "sign"
}

// User 积分
type User struct {
	Id             int64  `json:"id"`
	UserName       string `json:"userName"`
	PasswordDigest string `json:"passwordDigest"`
	Phone          string `json:"phone"`
	Integral       int    `json:"integral"`
	Location       string `json:"location"`
	Flg            int    `json:"flg"`
	Follower       int    `json:"follower"`
	Attention      int    `json:"attention"`
}

func (User) TableName() string {
	return "user"
}
