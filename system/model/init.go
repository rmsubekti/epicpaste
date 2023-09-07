package model

import (
	database "epicpaste/system/db"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	conn, sql := database.Connect()
	db = conn

	// init extensions, schemas, and functions
	sql.ExecFile("before_create_table", db)

	db.AutoMigrate(
		&Account{},
		&AccountSetting{},
		&User{},
		&Tag{},
		&Category{},
		&Syntax{},
		&Paste{},
	)

	// init table function and trigger relations
	sql.ExecFile("after_create_table", db)

	// init dummy data
	sql.ExecFile("dummy_user", db)
	sql.ExecFile("dummy_paste", db)

}
