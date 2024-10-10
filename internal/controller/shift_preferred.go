package controller

import (
	"github.com/gin-gonic/gin"

	"shiftmanager/internal/core/jwt"
	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/service"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/request"
)

type ShiftPreferredController struct {
	shiftPreferredService service.ShiftPreferredService
}

func NewShiftPreferredController() *ShiftPreferredController {
	return &ShiftPreferredController{
		shiftPreferredService: service.NewShiftPreferredService(),
	}
}

// GET /shift_preferred/:year/:month
func (ctr *ShiftPreferredController) ShiftPreferredPage(c *gin.Context) {
	c.HTML(200, "shift_preferred.html", gin.H{
		"year": c.Param("year"),
		"month": c.Param("month"),
	})
}

// GET /api/shift_preferred/:year/:month
func (ctr *ShiftPreferredController) ApiGetOne(c *gin.Context) {
	pl := jwt.GetPayload(c)

	var params request.ShiftPreferredUri
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

	c.JSON(200, result)
}

// POST /api/shift_preferred/:year/:month
func (ctr *ShiftPreferredController) ApiPost(c *gin.Context) {
	pl := jwt.GetPayload(c)

	var params request.ShiftPreferredUri
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


// GET /api/shift_preferred?year=:year&month=:month
func (ctr *ShiftPreferredController) ApiGet(c *gin.Context) {
	var params request.ShiftPreferredQuery
    if err := c.ShouldBindQuery(&params); err != nil {
        JsonError(c, 400, "不正なリクエストです。")
        return
	}

	var input dto.GetShiftPreferred
	utils.MapFields(&input, params)

	result, err := ctr.shiftPreferredService.Get(input)
	if err != nil {
		JsonError(c, 500, "シフト希望の取得に失敗しました。")
		return
	}

	c.JSON(200, result)
}