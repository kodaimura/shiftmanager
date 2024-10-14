package controller

import (
	"github.com/gin-gonic/gin"

	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/service"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/request"
)

type ShiftController struct {
	shiftService service.ShiftService
}

func NewShiftController() *ShiftController {
	return &ShiftController{
		shiftService: service.NewShiftService(),
	}
}

// GET /shifts/:year/:month
func (ctr *ShiftController) ShiftPage(c *gin.Context) {
	c.HTML(200, "shift.html", gin.H{
		"year": c.Param("year"),
		"month": c.Param("month"),
	})
}

// GET /api/shifts/:year/:month
func (ctr *ShiftController) ApiGetOne(c *gin.Context) {
	var params request.ShiftUri
    if err := c.ShouldBindUri(&params); err != nil {
        JsonError(c, 400, "不正なリクエストです。")
        return
	}

	var input dto.ShiftPK
	utils.MapFields(&input, params)

	result, err := ctr.shiftService.GetOne(input)
	if err != nil {
		JsonError(c, 500, "シフト希望の取得に失敗しました。")
		return
	}

	c.JSON(200, result)
}

// POST /api/shifts/:year/:month
func (ctr *ShiftController) ApiPost(c *gin.Context) {
	var params request.ShiftUri
    if err := c.ShouldBindUri(&params); err != nil {
        JsonError(c, 400, "不正なリクエストです。")
        return
	}

	var req request.PostShift
	if err := c.ShouldBindJSON(&req); err != nil {
		JsonError(c, 400, "不正なリクエストです。")
		return
	}

	var input dto.SaveShift
	utils.MapFields(&input, params)
	utils.MapFields(&input, req)

	if err := ctr.shiftService.Save(input); err != nil {
		JsonError(c, 500, "登録に失敗しました。")
		return
	}
	c.JSON(200, gin.H{})
}