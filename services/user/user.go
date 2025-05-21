package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/repositories"
)

func CreateSession(username, password string, db *repositories.Database, j *jwt.JwtConfig) (string, error) {

	u, err := db.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	token, err := j.LoginCheck(username, u.Password, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func AddUser(username, password, callerUsername string, admin bool, db *repositories.Database) error {
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

func DeleteUser(username, callerUsername string, db *repositories.Database) error {
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

func GetUser(username, callerUsername string, db *repositories.Database) (*gin.H, error) {
	u, err := db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil, err
	}

	if !u.Admin && username != callerUsername {
		return nil, errors.New("you lack permissions to get user(s)")
	}

	u, err = db.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	userProjects, err := db.GetProjectsByOwner(username)
	if err != nil {
		return nil, err
	}

	responseUserProjects := []ProjectResponse{}
	for _, p := range userProjects {
		responseUserProjects = append(responseUserProjects, ProjectResponse{p.Name, p.Owner, p.CreatedAt})
	}

	user := &gin.H{
		"username":  u.Username,
		"createdAt": u.CreatedAt,
		"projects":  responseUserProjects,
	}
	return user, nil
}

func GetUsers(callerUsername string, db *repositories.Database) (*gin.H, error) {
	u, err := db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil, err
	}
	if !u.Admin {
		return nil, errors.New("you lack permissions to get user(s)")
	}

	dbUsers, err := db.GetUsers()
	if err != nil {
		return nil, err
	}
	responseUsers := []UserResponse{}
	for _, u := range dbUsers {
		responseUsers = append(responseUsers, UserResponse{u.Username, u.Admin, u.CreatedAt})
	}

	users := &gin.H{
		"users": responseUsers,
	}
	return users, nil
}
