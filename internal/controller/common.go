package controller

import (
    "github.com/gin-gonic/gin"
)

func ResponseError (c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{"error": message})
    c.Abort()
}