package model

type User struct {
	ID   string `gorm:"type:varchar(60);primarykey:true;not null;unique" json:"id"`
	Name string `gorm:"type:varchar(40)" json:"name"`
}

func (User) TableName() string {
	return "user.user"
}

func (u *User) Get(id string) error {
	return db.First(&u, "id = ?", id).Error
}

func (u *User) Update(id string) error {
	temp := &User{}
	*temp = *u
	if err := u.Get(id); err != nil {
		return err
	}
	return db.Model(&u).Updates(&temp).Error
}
