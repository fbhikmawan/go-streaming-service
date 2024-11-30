package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {

	Rawtoken := c.GetHeader("Authorization")
	token := strings.Split(Rawtoken, "Bearer ")[1]

	if token == "" {
		c.JSON(401, gin.H{"error": "Authorization token not provided"})
		c.Abort()
		return
	}

	if token != "token" {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	

	c.Next()
}