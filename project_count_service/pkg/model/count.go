package model

type Fuzzy struct {
	Id   int64
	Type int
	Num  int
	Cid  int64
}

func (Fuzzy) TableName() string {
	return "fuzzy"
}

type Precise struct {
	Id   int64
	Type int
	Num  int
	Cid  int64
}

func (Precise) TableName() string {
	return "precise"
}
