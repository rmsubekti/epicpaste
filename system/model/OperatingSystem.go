package model

type OperatingSystem struct {
	ID   uint
	Name string
}

func (OperatingSystem) TableName() string {
	return "master.opeating_system"
}
