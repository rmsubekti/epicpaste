package db

import (
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//go:embed sql
var sqls embed.FS

type sqlEmbed struct {
	files embed.FS
}

func (e *sqlEmbed) ExecFile(fileName string, db *gorm.DB) {
	if sql, err := e.files.ReadFile("sql/" + fileName + ".sql"); err == nil {
		db.Exec(string(sql))
	} else {
		log.Fatal(err)
	}
}

func Connect() (*gorm.DB, *sqlEmbed) {
	var db *gorm.DB
	var err error
	if e := godotenv.Load(); e != nil {
		log.Println(e)
	}

	//database info
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOSTNAME")
	dbport := os.Getenv("POSTGRES_PORT")
	sslMode := os.Getenv("POSTGRES_SSLMODE")

	logLevel := logger.Silent
	if os.Getenv("EPIC_DEBUG") == "debug" {
		logLevel = logger.Info
	}

	//database uri
	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s ", host, user, password, dbName, dbport, sslMode)

	//connect to database
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(
			postgres.Open(dbUri),
			&gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
				Logger: logger.Default.LogMode(logLevel),
			},
		)
		if err == nil {
			log.Println("Connected to database host.")
			break
		}
		log.Println("\nReconnecting to your database host....")
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("No database found : %s", err)
	}

	return db, &sqlEmbed{files: sqls}
}
