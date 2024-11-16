package middleware

import (
	"net/http"
	"tender-backend/config"
	"tender-backend/internal/http/token"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		accessToken := authHeader

		claims, err := token.ExtractClaim(config.SecretKey, accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims", "details": err.Error()})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
