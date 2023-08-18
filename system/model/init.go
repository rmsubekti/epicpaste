package model

import (
	"epicpaste/system/config"
	"epicpaste/system/sql"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = config.GetDB()
	sql := sql.Populate()

	// init extensions, schemas, and functions
	sql.ExecFile("before_create_table", db)

	db.AutoMigrate(
		&Account{},
		&AccountSetting{},
		&User{},
		&Tag{},
		&Category{},
		&Paste{},
	)

	// init table function and trigger relations
	sql.ExecFile("after_create_table", db)

	// init dummy data
	sql.ExecFile("dummy_user", db)
	sql.ExecFile("dummy_paste", db)

}
