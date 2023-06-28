package model

type Language struct {
	ID            uint
	Name          string `gorm:"type:varchar(40)"`
	Description   string
	FileExtension string `gorm:"type:varchar(30)"`
}

func (Language) TableName() string {
	return "master.language"
}
