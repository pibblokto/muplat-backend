package models

import (
	"log"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

var db *gorm.DB

func ConnectDatabase() {
	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=testuser password=testpass dbname=testdb port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database")
	}

	db.AutoMigrate(&User{})
}
