package model

type Tag struct {
	ID   uint   `swaggerignore:"true"`
	Name string `gorm:"type:varchar(30)"`
}

func (Tag) TableName() string {
	return "master.tag"
}
