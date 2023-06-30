package utils

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func ExecSQLFile(fileName string, db *gorm.DB) {
	filePath, _ := filepath.Abs("./system/sql/" + fileName)
	if sql, err := os.ReadFile(filePath); err == nil {
		db.Exec(string(sql))
	} else {
		log.Fatal(err)
	}
}
