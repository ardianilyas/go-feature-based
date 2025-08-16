package middlewares

import (
	"net/http"
	"os"

	"github.com/ardianilyas/go-feature-based/pkg/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token, []byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_role", claims.Role)
		c.Set("user_id", claims.ID.String())
		c.Next()
	}
}