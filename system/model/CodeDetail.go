package model

type CodeDetail struct {
	ID          string `json:"-"`
	Description string
	Language    *Language          `gorm:"foreignKey:LangId" json:"lang,omitempty"`
	LangId      *uint              `json:"-"`
	App         *Application       `gorm:"foreignKey:AppId" json:"app,omitempty"`
	AppId       *uint              `json:"-"`
	OS          *[]OperatingSystem `gorm:"Many2Many:os_code;FOREIGNKEY:id;ASSOCIATION_FOREIGNKEY:id;" json:"os,omitempty"`
}

func (CodeDetail) TableName() string {
	return "master.code_detail"
}
