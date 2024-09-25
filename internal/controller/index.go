package controller

import (
	"github.com/gin-gonic/gin"
	
	"shiftmanager/internal/core/jwt"
)


type IndexController struct {}


func NewIndexController() *IndexController {
	return &IndexController{}
}


//GET /
func (ctr *IndexController) IndexPage(c *gin.Context) {
	pl := jwt.GetPayload(c)
	name := pl.AccountName

	c.HTML(200, "index.html", gin.H{
		"account_name": name,
	})
}