package controller

import (
	"github.com/gin-gonic/gin"

	"shiftmanager/internal/core/utils"
)

func JsonError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
	c.Abort()
}

func validateDayCSV(c *gin.Context, year, month int, value *string) bool {
	days, err := utils.ParseDayCSV(value)
	if err != nil {
		JsonError(c, 400, "不正なリクエストです。")
		return false
	}
	if err := utils.ValidateDaysInMonth(year, month, days); err != nil {
		JsonError(c, 400, "不正なリクエストです。")
		return false
	}
	return true
}
