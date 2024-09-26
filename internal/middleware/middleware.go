package middleware

import (
	"net/http"
	
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


func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := "admin"
		password := "pass"

		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != username || pass != password {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}