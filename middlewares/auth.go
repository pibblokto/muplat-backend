package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/pkg/jwt"
)

func JwtAuth(j *jwt.JwtConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !j.TokenValid(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid bearer token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
