package model

import (
	"epicpaste/system/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func createDummy(db *gorm.DB) {
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("5uperSecret"), bcrypt.DefaultCost)
	db.FirstOrCreate(&Account{
		ID:       "2539b6ba-7ff6-4864-b7bc-6fed752ee925",
		UserName: "admin",
		Email:    "su@bekti.net",
		Password: string(adminPassword),
	})
	db.FirstOrCreate(&User{
		ID:   "2539b6ba-7ff6-4864-b7bc-6fed752ee925",
		Name: utils.GenerateNewName(),
	})
}
