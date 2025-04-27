package repositories

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/muplat/muplat-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Connection       *gorm.DB
	InitUser         string `env:"INIT_USER" envDefault:"admin"`
	postgresHost     string `env:"POSTGRES_HOST"`
	postgresUser     string `env:"POSTGRES_USER"`
	postgresPassword string `env:"POSTGRES_PASSWORD"`
	postgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	database         string `env:"DATABASE"`
	initUserPassword string `env:"INIT_USER_PASSWORD"`
}

func NewDatabase() (db *Database) {

	err := env.Parse(db)
	if err != nil {
		log.Fatalf("Database config initialization error: %v", err)
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		db.postgresHost,
		db.postgresUser,
		db.postgresPassword,
		db.database,
		db.postgresPort,
	)
	db.Connection, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.Connection.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Deployment{},
		&models.AppConfig{},
		&models.PostgresConfig{},
	)
	err = db.initAdminUser()
	if err != nil {
		log.Fatalf("Failed to created init user: %v", err)
	}
	return db
}

func (db *Database) initAdminUser() error {
	var u models.User = models.User{
		Username: db.InitUser,
		Password: db.initUserPassword,
		Admin:    true,
	}
	err := db.Connection.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}
