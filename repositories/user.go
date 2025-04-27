package repositories

import (
	"github.com/muplat/muplat-backend/models"
)

func (db *Database) SaveUser(username, password string, admin bool) error {
	u := &models.User{
		Username: username,
		Password: password,
		Admin:    admin,
	}
	err := db.Connection.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetUserByUsername(username string) (*models.User, error) {
	u := &models.User{}
	err := db.Connection.Model(&models.User{}).Where("username = ?", username).Take(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (db *Database) DeleteUser(username string) error {
	u := &models.User{
		Username: username,
	}
	err := db.Connection.Model(&models.User{}).Delete(u).Error
	if err != nil {
		return err
	}
	return nil
}
