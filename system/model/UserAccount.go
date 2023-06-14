package model

import (
	"epicpaste/system/utils"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

type Account struct {
	ID       string `gorm:"type:varchar(40);primarykey:true;not null;unique"`
	UserName string `gorm:"type:varchar(60);not null;unique"`
	User     User   `gorm:"foreignKey:ID;references:ID"`
	Email    string `gorm:"type:varchar(40);not null;unique"`
	Password string `gorm:"type:char(200);not null"`
}

func (Account) TableName() string {
	return "user.account"
}

func (a *Account) Register() (err error) {

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

	if err = db.Create(&a).Error; err != nil {
		return
	}

	// create user profile
	return db.Create(&User{
		ID:   a.ID,
		Name: utils.GenerateNewName(),
	}).Error
}

func (a *Account) Login() error {
	password := a.Password

	if len(a.Email) < 1 && len(a.UserName) < 1 {
		return errors.New("email or username cannot be empty")
	}

	if result := db.Preload(clause.Associations).First(&a, "user_name = ?", a.UserName); result.RowsAffected < 1 {
		if result := db.Preload(clause.Associations).First(&a, "email = ?", a.Email); result.RowsAffected < 1 {
			return errors.New("email or username not registered")
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		return errors.New("wrong password")
	}
	return nil
}
