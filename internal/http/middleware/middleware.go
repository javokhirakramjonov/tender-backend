package middleware

import (
	"net/http"
	"tender-backend/internal/http/token"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		claims, err := token.ExtractClaim(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func ClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "client" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only clients can access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ContractorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "contractor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only contractors can access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}
