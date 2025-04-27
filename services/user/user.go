package user

import (
	"errors"

	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/repositories"
)

func CreateSession(username, password string, j *jwt.JwtConfig) (string, error) {
	token, err := j.LoginCheck(username, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func AddUser(username, password, callerUsername string, admin bool, db *repositories.Database, j *jwt.JwtConfig) error {
	u, err := db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil
	}

	if !u.Admin {
		return errors.New("you lack permissions to add user")
	}

	err = db.SaveUser(username, password, admin)
	if err != nil {
		return nil
	}
	return nil
}

func DeleteUser(username, callerUsername string, db *repositories.Database, j *jwt.JwtConfig) error {
	u, err := db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil
	}

	if !u.Admin {
		return errors.New("you lack permissions to delete user")
	}

	if username == db.InitUser {
		return errors.New("init user can't be deleted")
	}

	err = db.DeleteUser(username)
	if err != nil {
		return nil
	}
	return nil
}
