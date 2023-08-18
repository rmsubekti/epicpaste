package sql

import (
	"embed"
	"log"

	"gorm.io/gorm"
)

//go:embed *.sql
var sqls embed.FS

type sqlEmbed struct {
	files embed.FS
}

func Populate() *sqlEmbed {
	return &sqlEmbed{files: sqls}
}
func (e *sqlEmbed) ExecFile(fileName string, db *gorm.DB) {
	if sql, err := e.files.ReadFile(fileName + ".sql"); err == nil {
		db.Exec(string(sql))
	} else {
		log.Fatal(err)
	}
}
