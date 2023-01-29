package model

type Follower struct {
	Id         int64
	UserId     int64
	FollowerId int64
}

func (Follower) TableName() string {
	return "follower"
}
