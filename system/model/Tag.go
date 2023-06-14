package model

type Tag struct {
	ID   uint
	Name string
}

func (Tag) TableName() string {
	return "master.tag"
}
