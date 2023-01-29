package model

type Comment struct {
	Id       int64
	Type     int
	MemberId int64
	State    int
	Content  string
}

func (Comment) TableName() string {
	return "comment"
}
