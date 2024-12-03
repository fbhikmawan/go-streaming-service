package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

var authService = services.NewAuthService()

func AuthMiddleware(c *gin.Context) {

	Rawtoken := c.GetHeader("Authorization")
	token := strings.Split(Rawtoken, "Bearer ")[1]

	if token == "" {
		c.JSON(401, gin.H{"error": "Authorization token not provided"})
		c.Abort()
		return
	}

	_, err := authService.ValidateToken(token)

	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Next()
}