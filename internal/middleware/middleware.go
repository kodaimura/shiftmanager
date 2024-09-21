package middleware

import (
	"github.com/gin-gonic/gin"
	"shiftmanager/internal/core/jwt"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwt.Auth(c); err != nil {
			c.Redirect(303, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}


func JwtAuthApiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwt.Auth(c); err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}