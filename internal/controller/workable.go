package controller

import (
    "log"
    "sort"
    "strings"

    "github.com/gin-gonic/gin"
    
    "shiftmanager/pkg/utils"
    "shiftmanager/internal/pkg/jwtauth"
    "shiftmanager/internal/model/repository"
    "shiftmanager/internal/model/entity"
    "shiftmanager/internal/constants"
)


type WorkableController interface {
    WorkablePage(c *gin.Context)
    Workable(c *gin.Context)
}


type workableController struct {
    ur repository.UserRepository
    wr repository.WorkableRepository
}


func NewWorkableController() WorkableController {
    ur := repository.NewUserRepository()
    wr := repository.NewWorkableRepository()
    return workableController{ur, wr}
}


//GET /workable/:year/:month
func (wc workableController)WorkablePage(c *gin.Context) {
    username, err := jwtauth.ExtractUserName(c)
    uid, _ := jwtauth.ExtractUId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    year := c.Param("year")
    month := c.Param("month")

    workable, _ := wc.wr.SelectByUIdYearMonth(uid, year, month)

    c.HTML(200, "workable.html", gin.H{
        "appname": constants.AppName,
        "username": username,
        "workable": workable,
        "year": year,
        "month": month,
    })
}


//POST /workable/:year/:month
func (wc workableController)Workable(c *gin.Context) {
    uid, err := jwtauth.ExtractUId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    w := &entity.Workable{}
    w.UId = uid
    w.Year = c.Param("year")
    w.Month = c.Param("month")

    wd, err := sortDays(c.PostForm("workable"))

    if err != nil {
        c.Redirect(303, "/login")
        return
    }
    w.WorkableDays = wd

    if wc.wr.Upsert(*w) != nil {
        log.Panic(err)
        return
    }

    c.Redirect(303, "/")
}


func sortDays(csvDays string) (string, error){
    if csvDays == "" {
        return "", nil
    }
    days := strings.Split(csvDays, ",")
    
    isl, err := utils.AtoiSlice(days)

    if err != nil {
        return "", err
    }

    sort.Ints(isl)
    return strings.Join(utils.ItoaSlice(isl), ","), nil
}
