package model

type Syntax struct {
	ID   uint   `swaggerignore:"true"`
	Name string `gorm:"type:varchar(60)"`
}

type Syntaxs []Syntax

func (Syntax) TableName() string {
	return "master.syntax"
}

func (s *Syntaxs) List() error {
	return db.Find(&s).Error
}

func (s *Syntax) Get(name string) error {
	return db.First(&s, "name=?", s.Name).Error
}
func (s *Syntax) SetAvailable() error {
	if s != nil && s.ID == 0 {
		return s.Get(s.Name)
	}
	return nil
}
