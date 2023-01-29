package model

type Attention struct {
	Id          int64
	UserId      int64
	AttentionId int64
	Flg         int
}

func (Attention) TableName() string {
	return "attention"
}
