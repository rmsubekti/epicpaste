package model

type OperatingSystem struct {
	ID   uint
	Name string `gorm:"type:varchar(90)"`
}

func (OperatingSystem) TableName() string {
	return "master.opeating_system"
}
