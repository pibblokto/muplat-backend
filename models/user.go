package models

import (
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
	Projects  []Project `gorm:"foreignKey:Owner;references:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *User) SaveUser() (*User, error) {
	err := db.Create(u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
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
