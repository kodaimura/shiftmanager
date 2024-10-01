package controller

import (
	"github.com/gin-gonic/gin"

	"shiftmanager/internal/core/jwt"
	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/service"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/request"
	"shiftmanager/internal/response"
)

type ShiftPreferredController struct {
	shiftPreferredService service.ShiftPreferredService
}

func NewShiftPreferredController() *ShiftPreferredController {
	return &ShiftPreferredController{
		shiftPreferredService: service.NewShiftPreferredService(),
	}
}

// GET /shift_preferred
func (ctr *ShiftPreferredController) ShiftPreferredPage(c *gin.Context) {
	c.HTML(200, "shift_preferred.html", gin.H{})
}

// GET /api/shift_preferred/:year/:month
func (ctr *ShiftPreferredController) ApiGetOne(c *gin.Context) {
	pl := jwt.GetPayload(c)

	var params request.GetShiftPreferredParams
    if err := c.ShouldBindUri(&params); err != nil {
        JsonError(c, 400, "不正なリクエストです。")
        return
	}

	var input dto.ShiftPreferredPK
	utils.MapFields(&input, params)
	input.AccountId = pl.AccountId

	result, err := ctr.shiftPreferredService.GetOne(input)
	if err != nil {
		JsonError(c, 500, "シフト希望の取得に失敗しました。")
		return
	}

	var res response.GetShiftPreferred
	utils.MapFields(&res, result)
	c.JSON(200, res)
}

// POST /api/shift_preferred/:year/:month
func (ctr *ShiftPreferredController) ApiPost(c *gin.Context) {
	pl := jwt.GetPayload(c)

	var params request.PostShiftPreferredParams
    if err := c.ShouldBindUri(&params); err != nil {
        JsonError(c, 400, "不正なリクエストです。")
        return
	}

	var req request.PostShiftPreferred
	if err := c.ShouldBindJSON(&req); err != nil {
		JsonError(c, 400, "不正なリクエストです。")
		return
	}

	var input dto.SaveShiftPreferred
	utils.MapFields(&input, params)
	utils.MapFields(&input, req)
	input.AccountId = pl.AccountId

	if err := ctr.shiftPreferredService.Save(input); err != nil {
		JsonError(c, 500, "登録に失敗しました。")
		return
	}
	c.JSON(200, gin.H{})
}