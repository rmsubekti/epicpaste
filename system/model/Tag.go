package model

type Tag struct {
	ID   uint
	Name string `gorm:"type:varchar(30)"`
}

func (Tag) TableName() string {
	return "master.tag"
}
