package model

type AccountSetting struct {
	ID        string `gorm:"type:varchar(40);primarykey:true;not null;unique"`
	Crawlable bool
}

func (AccountSetting) TableName() string {
	return "user.account_setting"
}

func (a *AccountSetting) Get(id string) error {
	return db.First(&a, "id = ?", id).Error
}

func (a *AccountSetting) Update(id string) error {
	temp := &AccountSetting{}
	*temp = *a
	if err := a.Get(id); err != nil {
		return err
	}
	return db.Model(&a).Updates(&temp).Error
}
