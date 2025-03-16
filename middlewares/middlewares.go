package middlewares

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(">>> Middleware before code")
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
		fmt.Println(">>> Middleware after code")
	}
}
