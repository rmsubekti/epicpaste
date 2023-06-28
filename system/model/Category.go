package model

type Category struct {
	ID   uint
	Name string `gorm:"type:varchar(40)"`
}

func (Category) TableName() string {
	return "master.category"
}
