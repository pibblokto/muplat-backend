package models

import (
	"fmt"
	"log"

	"github.com/muplat/muplat-backend/pkg/setup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	cfg setup.MuplatCfg = setup.LoadConfig()
)

func ConnectDatabase() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.Database,
		cfg.PostgresPort,
	)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.AutoMigrate(
		&User{},
		&Project{},
		&Deployment{},
		&AppConfig{},
		&PostgresConfig{},
	)
}

func CreateInitUser() {
	var u User = User{
		Username: cfg.InitUser,
		Password: cfg.InitUserPassword,
	}
	err := db.Create(&u).Error
	if err != nil {
		log.Fatalf("Failed to create init user: %v", err)
	}
}
