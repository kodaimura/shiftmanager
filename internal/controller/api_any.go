package controller

import (
	"shiftmanager/internal/pkg/jwtauth"
    "shiftmanager/internal/model/repository"
    "shiftmanager/internal/model/entity"

    "github.com/gin-gonic/gin"
)


//GET /api/group/workables/:year/:month
func GroupWorkables(c *gin.Context) {
	gid, err := jwtauth.ExtractGId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

	year := c.Param("year")
    month := c.Param("month")

	wr := repository.NewWorkableRepository()
	workables, err :=  wr.GetWorkableExp1ByGId(gid, year, month)

    c.JSON(200, workables)
}


//POST /api/shift/:year/:month
func SaveShift(c *gin.Context) {
    gid, err := jwtauth.ExtractGId(c)
    var body map[string]interface{}
    c.BindJSON(&body)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }
    
    year := c.Param("year")
    month := c.Param("month")

    s := &entity.Shift{}
    s.Year = year
    s.Month = month
    s.StoreHoliday = body["storeholiday"].(string)
    s.Shift = body["shift"].(string)
    s.GId = gid

    sr := repository.NewShiftRepository()
    err = sr.Upsert(*s)

    if err != nil {
        c.JSON(500, gin.H{})
        return
    }

    c.JSON(200, gin.H{})
}
