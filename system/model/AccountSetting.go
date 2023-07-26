package model

type AccountSetting struct {
	ID        string `gorm:"type:varchar(40);primarykey:true;not null;unique"`
	Crawlable bool
}

func (AccountSetting) TableName() string {
	return "user.account_setting"
}
