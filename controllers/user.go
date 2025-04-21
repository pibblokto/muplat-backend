package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/models"
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/pkg/setup"
	"golang.org/x/crypto/bcrypt"
)

var (
	InitUser string = setup.LoadConfig().InitUser
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Admin    *bool  `json:"admin" binding:"required"`
}

type UserDeleteInput struct {
	Username string `json:"username" binding:"required"`
}

func LoginCheck(username string, password string) (string, error) {
	u := models.User{}

	err := u.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(username)

	if err != nil {
		return "", err
	}

	return token, nil
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := LoginCheck(input.Username, input.Password)
	fmt.Println(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func AddUser(c *gin.Context) {
	var input UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = models.IsUserAdmin(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	u.Username = input.Username
	u.Password = input.Password
	u.Admin = *input.Admin

	_, err = u.SaveUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("added user %s", input.Username)})
}

func DeleteUser(c *gin.Context) {
	var input UserDeleteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	err = u.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !u.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("you have to be an admin in order to delete users")})
		return
	}
	if input.Username == InitUser {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("can't delete initial user")})
		return
	}

	u = models.User{}
	err = u.GetUserByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = u.DeleteUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("deleted user %s", input.Username)})
}
