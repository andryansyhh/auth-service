package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/andryansyhh/auth-service/pkg/domain/dto"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		switch {
		case errors.Is(err.Err, dto.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Err.Error()})
		case errors.Is(err.Err, dto.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Err.Error()})
		case errors.Is(err.Err, dto.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Err.Error()})
		default:
			log.Printf("Internal server error: %v", err.Err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		}
	}
}
