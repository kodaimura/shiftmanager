package controller

import (
	"github.com/gin-gonic/gin"

	"shiftmanager/internal/core/jwt"
	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/service"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/request"
)

type AccountProfileController struct {
	accountProfileService service.AccountProfileService
}

func NewAccountProfileController() *AccountProfileController {
	return &AccountProfileController{
		accountProfileService: service.NewAccountProfileService(),
	}
}

// GET /account_profile
func (ctr *AccountProfileController) AccountProfilePage(c *gin.Context) {
	c.HTML(200, "account_profile.html", gin.H{})
}

// GET /api/account_profile
func (ctr *AccountProfileController) ApiGetOne(c *gin.Context) {
	pl := jwt.GetPayload(c)
	result, err := ctr.accountProfileService.GetOne(pl.AccountId)

	if err != nil {
		JsonError(c, 500, "プロフィール情報の取得に失敗しました。")
		return
	}

	c.JSON(200, result)
}

// POST /api/account_profile
func (ctr *AccountProfileController) ApiPost(c *gin.Context) {
	pl := jwt.GetPayload(c)

	var req request.PostAccountProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		JsonError(c, 400, "不正なリクエストです。")
		return
	}

	var input dto.SaveAccountProfile
	utils.MapFields(&input, req)

	if err := ctr.accountProfileService.Save(pl.AccountId, input); err != nil {
		JsonError(c, 500, "登録に失敗しました。")
		return
	}
	c.JSON(200, gin.H{})
}