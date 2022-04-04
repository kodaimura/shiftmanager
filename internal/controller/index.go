package controller

import (
    "time"
    "strconv"
    "github.com/gin-gonic/gin"
    
    "shiftmanager/internal/pkg/jwtauth"
    "shiftmanager/internal/model/repository"
    "shiftmanager/internal/constants"
)


type IndexController interface {
    IndexPage(c *gin.Context)
}


type indexController struct {
    ur repository.UserRepository
    pr repository.ProfileRepository
    wr repository.WorkableRepository
}


func NewIndexController() IndexController {
    ur := repository.NewUserRepository()
    pr := repository.NewProfileRepository()
    wr := repository.NewWorkableRepository()
    return indexController{ur, pr, wr}
}


func (ic indexController)IndexPage(c *gin.Context) {
    username, err := jwtauth.ExtractUserName(c)
    uid, _ := jwtauth.ExtractUId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }
    now := time.Now()
    y, m, _ := now.Date()
    d := time.Date(y, m+1, 1, 0, 0, 0, 0, now.Location())

    profile, _ := ic.pr.GetProfileExp1ByUId(uid)
    workabledays, _ := ic.wr.GetWorkableExp1ByGId(profile.GId, strconv.Itoa(d.Year()), strconv.Itoa(int(d.Month())))
    

    c.HTML(200, "index.html", gin.H{
        "appname": constants.AppName,
        "username": username,
        "profile": profile,
        "workabledays": workabledays,
    })
}

