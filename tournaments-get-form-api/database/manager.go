package database

import (
	"fmt"
	"log"
	"os"
	"tournaments-api/classes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBase struct {
	DataBase *gorm.DB
}

func Init() DataBase {
	var DBManager DataBase

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, "disable",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}

	db.AutoMigrate(&classes.Tournament{}, &classes.Result{}, &classes.Metric{}, &classes.Metadata{})

	DBManager.DataBase = db

	return DBManager
}
