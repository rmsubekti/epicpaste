package model

type Category struct {
	ID   uint
	Name string
}

func (Category) TableName() string {
	return "master.category"
}
