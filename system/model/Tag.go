package model

type Tag struct {
	ID   uint   `swaggerignore:"true"`
	Name string `gorm:"type:varchar(30)"`
}

type Tags []Tag

func (Tag) TableName() string {
	return "master.tag"
}

func (t *Tags) List() error {
	return db.Find(&t).Error
}

func (t *Tag) Get(name string) error {
	return db.First(&t, "name = ?", name).Error
}

func (t *Tags) SetAvailable() (err error) {
	for i, v := range *t {
		err = v.Get(v.Name)
		(*t)[i] = v
	}
	return err
}
