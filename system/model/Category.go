package model

type Category struct {
	ID   uint   `swaggerignore:"true"`
	Name string `gorm:"type:varchar(40)"`
}

type Categories []Category

func (Category) TableName() string {
	return "master.category"
}

func (c *Categories) List() error {
	return db.Find(&c).Error
}

func (c *Category) Get(name string) error {
	return db.First(&c, "name=?", c.Name).Error
}
func (c *Category) SetAvailable() error {
	if c != nil && c.ID == 0 {
		return c.Get(c.Name)
	}
	return nil
}
