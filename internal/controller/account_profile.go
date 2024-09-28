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
	pl := jwt.GetPayload(c)
	profile, err := ctr.accountProfileService.GetOne(pl.AccountId)

	if err != nil {
		c.HTML(200, "account_profile.html", gin.H{
			"display_name": "",
			"account_role": "",
		})
		c.Abort()
		return
	}

	c.HTML(200, "account_profile.html", gin.H{
		"display_name": profile.DisplayName,
		"account_role": profile.AccountRole,
	})
}

// POST /api/account_profile
func (ctr *AccountProfileController) PostAccountProfile(c *gin.Context) {
	pl := jwt.GetPayload(c)

	var req request.PostAccountProfile
	if err := c.ShouldBind(&req); err != nil {
		JsonError(c, 400, "不正なリクエストです。")
		return
	}

	var input dto.SaveAccountProfile
	utils.MapFields(&input, req)

	if err := ctr.accountProfileService.Save(pl.AccountId, input); err != nil {
		c.HTML(200, "account_profile.html", gin.H{
			"display_name": input.DisplayName,
			"account_role": input.AccountRole,
			"error": "登録に失敗しました。",
		})
		c.Abort()
		return
	}

	c.Redirect(303, "/")
}