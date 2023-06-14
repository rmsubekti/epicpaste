package model

type Application struct {
	ID   uint
	Name string
}

func (Application) TableName() string {
	return "master.application"
}
