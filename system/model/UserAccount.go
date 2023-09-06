package model

import (
	"epicpaste/system/utils"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

type Account struct {
	ID string `json:"id" gorm:"type:varchar(40);primarykey:true;not null;unique" swaggerignore:"true"`
	//any combination of alphanumeric charachter allowed
	UserName string `json:"username" gorm:"type:varchar(60);not null;unique"`
	//must be valid email address
	//cannot start with or end with non alphanumeric
	//charachter (`.` and `-`) are allowed
	Email string `json:"email" gorm:"type:varchar(40);not null;unique"`
	//any string character
	Password string         `json:"password" gorm:"type:char(200);not null"`
	User     User           `json:"user" gorm:"foreignKey:UserName;references:UserName" swaggerignore:"true"`
	Setting  AccountSetting `json:"setting" gorm:"foreignKey:ID;references:ID" swaggerignore:"true"`
}

var (
	mailRegex   = regexp.MustCompile(`^(?:[[:alnum:]]+[[:alnum:]\-\.]+[[:alnum:]])+@(?:[[:alnum:]]+[[:alnum:]\-\.]+[[:alnum:]])+\.(?:[[:alpha:]]{2,6})$`)
	allNumRegex = regexp.MustCompile(`^(?:[[:alnum:]]+)$`)
)

func (Account) TableName() string {
	return "user.account"
}

type ChangePassword struct {
	Current string `json:"current"`
	New     string `json:"new"`
	Confirm string `json:"confirm"`
}

func (c *ChangePassword) validate() (err error) {
	if c.Current == "" {
		return fmt.Errorf("current password cannot be empty")
	}
	if c.New == "" {
		return fmt.Errorf("new password cannot be empty")
	}
	if c.Confirm == "" {
		return fmt.Errorf("confirm password cannot be empty")
	}
	if c.Current == c.New {
		return fmt.Errorf("current password is already used")
	}
	if c.New != c.Confirm {
		return fmt.Errorf("confirm password should be equal to new password")
	}

	return nil
}

func (a *Account) Register() (err error) {
	if !mailRegex.MatchString(a.Email) {
		return errors.New("email not valid")
	}

	if !allNumRegex.MatchString(a.UserName) {
		return errors.New("username should only contain alphanumeric")
	}

	account := &Account{}
	if result := db.First(account, "email = ?", a.Email); result.RowsAffected > 0 {
		return errors.New("email already registered")
	}

	if result := db.First(account, "user_name = ?", a.UserName); result.RowsAffected > 0 {
		return errors.New("username already used")
	}

	if len(a.Password) <= 8 {
		return errors.New("password at least have 8 or more characters")
	}
	a.ID = uuid.NewString()

	// save hashed password to db
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	a.Password = string(hashedPassword)
	a.User = User{
		UserName: a.UserName,
		Name:     utils.GenerateNewName(),
	}
	a.Setting = AccountSetting{
		ID:        a.ID,
		Crawlable: false,
	}

	if err = db.Create(&a).Error; err != nil {
		return
	}
	return
}

func (a *Account) Login() error {
	password := a.Password

	if len(a.UserName) < 1 {
		return errors.New("email or username cannot be empty")
	}

	if result := db.Preload(clause.Associations).First(&a, "user_name = ?", a.UserName); result.RowsAffected < 1 {
		if result := db.Preload(clause.Associations).First(&a, "email = ?", a.UserName); result.RowsAffected < 1 {
			return errors.New("email or username not registered")
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		return errors.New("wrong password")
	}
	return nil
}

func (a *Account) Get(username string) error {
	return db.Preload(clause.Associations).First(&a, "user_name = ?", username).Error
}

func (a *Account) Update(username string) error {
	temp := &Account{}
	*temp = *a
	temp.Password = ""

	if a.Email != "" && !mailRegex.MatchString(a.Email) {
		return errors.New("your email is not valid")
	}

	if a.UserName != "" && !allNumRegex.MatchString(a.UserName) {
		return errors.New("username should be using alphanumeric character")
	}

	if err := a.Get(username); err != nil {
		return err
	}
	return db.Model(&a).Updates(&temp).Error
}

func (a *Account) ChangePassword(password ChangePassword) (err error) {
	if err = password.validate(); err != nil {
		return
	}

	if err = a.Get(a.UserName); err != nil {
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password.New)); err != nil {
		return errors.New("wrong current password")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password.New), bcrypt.DefaultCost)
	password.New = string(hashedPassword)

	if err = db.Model(&a).Updates(&Account{Password: password.New}).Error; err != nil {
		return
	}

	return
}

func (a *Account) ChangeEmail(mail string) (err error) {

	if !mailRegex.MatchString(a.Email) {
		return errors.New("new email is not valid")
	}

	if a.Email == mail {
		return errors.New("new email is currently used")
	}

	if result := db.First(&Account{}, "email=?", mail); result.RowsAffected > 0 {
		return errors.New("new email is already registered")
	}
	if err = db.Model(&a).Updates(&Account{Email: mail}).Error; err != nil {
		return
	}
	return
}
