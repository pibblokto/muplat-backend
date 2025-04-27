package jwt

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (j *JwtConfig) GenerateToken(username string) (string, error) {
	lifespanMinutes, err := strconv.Atoi(j.JwtLifespanMinutes)
	jwtSecret := j.JwtSecret

	if err != nil {
		return "", nil
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(lifespanMinutes)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (j *JwtConfig) TokenValid(c *gin.Context) bool {
	jwtSecret := j.JwtSecret
	tokenString := j.ExtractToken(c)
	_, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		log.Printf("Token validation error: %v\n", err)
		return false
	}
	return true
}

func (j *JwtConfig) ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func (j *JwtConfig) ExtractTokenUsername(c *gin.Context) (string, error) {
	tokenString := j.ExtractToken(c)
	jwtSecret := j.JwtSecret

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to cast token.Claims to jwt.MapClaims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("failed to cast claims[\"username\"] to string")
	}
	return username, nil
}

func (j *JwtConfig) LoginCheck(username, hashedPassword, password string) (string, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := j.GenerateToken(username)

	if err != nil {
		return "", err
	}

	return token, nil
}
