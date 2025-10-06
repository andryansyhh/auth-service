package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		// Expect header: Authorization: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization format"})
			return
		}

		tokenStr := parts[1]
		claims, err := jwtManager.Verify(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Store user info in context for later use
		c.Set("user_id", claims.ID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
