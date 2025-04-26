package repositories

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/muplat/muplat-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Connection       *gorm.DB
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	Database         string `env:"DATABASE"`
	InitUser         string `env:"INIT_USER" envDefault:"admin"`
	InitUserPassword string `env:"INIT_USER_PASSWORD"`
}

func (db *DatabaseConfig) InitDatabase() {

	err := env.Parse(db)
	if err != nil {
		log.Fatalf("Database config initialization error: %v", err)
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		db.PostgresHost,
		db.PostgresUser,
		db.PostgresPassword,
		db.Database,
		db.PostgresPort,
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
}

func (db *DatabaseConfig) initAdminUser() error {
	var u models.User = models.User{
		Username: db.InitUser,
		Password: db.InitUserPassword,
		Admin:    true,
	}
	err := db.Connection.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}
