package controller

import (
    "log"
    "database/sql"
    "github.com/gin-gonic/gin"
    
    "shiftmanager/internal/pkg/jwtauth"
    "shiftmanager/internal/model/entity"
    "shiftmanager/internal/model/repository"
    "shiftmanager/internal/constants"
)


type ProfileController interface {
    ProfilePage(c *gin.Context)
    Profile(c *gin.Context)
}


type profileController struct {
    ur repository.UserRepository
    pr repository.ProfileRepository
    gr repository.GroupRepository
}


func NewProfileController() ProfileController {
    ur := repository.NewUserRepository()
    pr := repository.NewProfileRepository()
    gr := repository.NewGroupRepository()

    return profileController{ur, pr, gr}
}


func (pc profileController)ProfilePage(c *gin.Context) {
    username, err := jwtauth.ExtractUserName(c)
    uid, _ := jwtauth.ExtractUId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }
    
    gr := repository.NewGeneralRepository()
    roles, _ := gr.SelectByClass("role")
    profile, _ := pc.pr.GetProfileExp1ByUId(uid)

    c.HTML(200, "profile.html", gin.H{
        "appname": constants.AppName,
        "username": username,
        "roles": roles,
        "profile": profile,
    })
}


func (pc profileController)Profile(c *gin.Context) {
    uid, err := jwtauth.ExtractUId(c)
    username, _ := jwtauth.ExtractUserName(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    groupName := c.PostForm("groupname")
    group, err := pc.gr.SelectByGroupName(groupName)
    if err == sql.ErrNoRows {
        err = pc.gr.Insert(groupName)
        group, err = pc.gr.SelectByGroupName(groupName)  
    }

    if err != nil {
        log.Panic(err)
        return
    }
    p := &entity.Profile{} 
    p.UId = uid
    p.GId = group.GId
    p.Role = c.PostForm("role")
    p.AbbName = c.PostForm("abbname")

    err = pc.pr.Upsert(*p)

    if err != nil {
        c.Redirect(303, "/")
        return 
    }

    jwtString, _ := jwtauth.GenerateJWT(uid, username, group.GId)
    c.SetCookie(jwtauth.JwtKeyName, jwtString, constants.CookieExpires, "/", constants.HostName, false, true)

    c.Redirect(303, "/")
}
