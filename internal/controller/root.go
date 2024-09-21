package controller

import (
	"github.com/gin-gonic/gin"
	
	"shiftmanager/internal/core/jwt"
)


type RootController struct {}


func NewRootController() *RootController {
	return &RootController{}
}


//GET /
func (rc *RootController) IndexPage(c *gin.Context) {
	pl := jwt.GetPayload(c)
	name := pl.AccountName

	c.HTML(200, "index.html", gin.H{
		"account_name": name,
	})
}