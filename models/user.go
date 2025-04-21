package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username  string `gorm:"primarykey"`
	Password  string
	CreatedAt time.Time
	Admin     bool
	Projects  []Project `gorm:"foreignKey:Owner;references:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *User) SaveUser() (*User, error) {
	err := db.Create(u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) GetUserByUsername(username string) error {
	err := db.Model(&User{}).Where("username = ?", username).Take(u).Error
	if err != nil {
		return err
	}
	return nil
}

func IsUserAdmin(username string) error {
	u := &User{}
	err := db.Model(&User{}).Where("username = ?", username).Take(u).Error
	if err != nil {
		return err
	}
	if !u.Admin {
		return errors.New("caller is not an admin")
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

func (u *User) DeleteUser() error {
	err := db.Model(&User{}).Delete(u).Error
	if err != nil {
		return err
	}
	return nil
}
