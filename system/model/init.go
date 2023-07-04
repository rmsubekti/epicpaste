package model

import (
	"epicpaste/system/config"
	"epicpaste/system/utils"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = config.GetDB()

	// init extensions, schemas, and functions
	utils.ExecSQLFile("before_table.sql", db)

	db.AutoMigrate(
		&Account{},
		&User{},
		&Tag{},
		&Category{},
		&Paste{},
	)

	// init table function and trigger relations
	utils.ExecSQLFile("after_table.sql", db)

	// init dummy data
	utils.ExecSQLFile("dummy_user.sql", db)
	utils.ExecSQLFile("dummy_paste.sql", db)

}
