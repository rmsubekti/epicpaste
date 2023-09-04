package model

type User struct {
	UserName string `gorm:"type:varchar(60);primarykey:true;not null;unique" json:"username"`
	Name     string `gorm:"type:varchar(40)" json:"name"`
}

func (User) TableName() string {
	return "user.user"
}

func (u *User) Get(username string) error {
	return db.First(&u, "user_name = ?", username).Error
}

func (u *User) Update() (err error) {
	temp := &User{}
	*temp = *u

	if err = u.Get(temp.UserName); err != nil {
		return
	}

	if err = db.Model(&u).Updates(&temp).Error; err != nil {
		return
	}

	return
}
