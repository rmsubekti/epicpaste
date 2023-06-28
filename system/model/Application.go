package model

type Application struct {
	ID   uint
	Name string `gorm:"type:varchar(90)"`
}

func (Application) TableName() string {
	return "master.application"
}
