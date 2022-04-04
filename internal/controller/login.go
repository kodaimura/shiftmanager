package controller

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    
    "shiftmanager/internal/pkg/jwtauth"
    "shiftmanager/internal/dto"
    "shiftmanager/internal/model/repository"
    "shiftmanager/internal/constants"
)


type LoginController interface {
    LoginPage(c *gin.Context)
    Login(c *gin.Context)
    Logout(c *gin.Context)
}


type loginController struct {
    ur repository.UserRepository
}


func NewLoginController() LoginController {
    ur := repository.NewUserRepository()
    return loginController{ur}
}


func (lc loginController)LoginPage(c *gin.Context) {
    c.HTML(200, "login.html", gin.H{
        "appname": constants.AppName,
    })
}


func (lc loginController)Login(c *gin.Context) {
    ld := &dto.LoginDto{}
    ld.UserName = c.PostForm("username")
    ld.Password = c.PostForm("password")

    user, err := lc.ur.SelectByUserName(ld.UserName)

    if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ld.Password)) != nil{
        c.HTML(401, "login.html", gin.H{
            "appname":constants.AppName,
            "error": "UsernameまたはPasswordが異なります。",
        })
        c.Abort()
        return
    }

    profile, err :=  repository.NewProfileRepository().SelectByUId(user.UId)
    gid := -1

    if err == nil {
        gid = profile.GId
    } 

    jwtString, err := jwtauth.GenerateJWT(user.UId, user.UserName, gid)
    if err != nil {
        c.HTML(500, "login.html", gin.H{
            "appname": constants.AppName,
            "error": "ログインに失敗しました。",
        })
        c.Abort()
        return
    }

    c.SetCookie(jwtauth.JwtKeyName, jwtString, constants.CookieExpires, "/", constants.HostName, false, true)
    c.Redirect(303, "/")
}


func (lc loginController)Logout(c *gin.Context) {
    c.SetCookie(jwtauth.JwtKeyName, "", 0, "/", constants.HostName, false, true)
    c.Redirect(303, "/login")
}
