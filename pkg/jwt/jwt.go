package jwt

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/muplat/muplat-backend/pkg/setup"
)

var cfg setup.MuplatCfg = setup.LoadConfig()

func GenerateToken(username string) (string, error) {
	lifespanHours, err := strconv.Atoi(cfg.JwtLifespanHours)
	jwtSecret := cfg.JwtSecret

	if err != nil {
		return "", nil
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(lifespanHours)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func TokenValid(c *gin.Context) bool {
	jwtSecret := cfg.JwtSecret
	tokenString := ExtractToken(c)
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

func ExtractToken(c *gin.Context) string {
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

func ExtractTokenUsername(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	jwtSecret := cfg.JwtSecret

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
