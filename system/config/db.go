package config

import (
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

func GetDB() *gorm.DB {
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
	logLevel := logger.Default.LogMode(logger.Silent)
	if os.Getenv("DEBUG") == "debug" {
		logLevel = logger.Default.LogMode(logger.Info)
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
				Logger: logLevel,
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

	return db
}
