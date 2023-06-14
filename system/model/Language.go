package model

type Language struct {
	ID            uint
	Name          string
	Description   string
	FileExtension string
}

func (Language) TableName() string {
	return "master.language"
}
