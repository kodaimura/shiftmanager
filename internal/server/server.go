package server

import (
    "log"
    "io"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"

    "shiftmanager/internal/controller"
    "shiftmanager/internal/constants"
    "shiftmanager/internal/pkg/jwtauth"
)


func Run() {
    setLogger()
    r := router()
    r.Run(constants.Port)
}


func setLogger () {
    logfolder := "log"
    logfile := "log/app.log"

    if _, err := os.Stat(logfolder); err != nil {
        os.Mkdir(logfolder, 0777)
    }

    if _, err := os.Stat(logfile); err == nil {
        t := time.Now()
        format := "2006-01-02-15-04-05"
        fname := "log/~" + t.Format(format) + ".log"
        if err := os.Rename(logfile, fname); err != nil {
            log.Panic(err)
        }
    }

    f, err := os.Create(logfile); 

    if err != nil {
        log.Panic(err)
    }

    gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
}


func router() *gin.Engine {
    r := gin.Default()
    
    //TEMPLATE
    r.LoadHTMLGlob("web/template/*.html")

    //STATIC
    r.Static("/css", "web/static/css")
    r.Static("/js", "web/static/js")

    
    //ROOT
    lc := controller.NewLoginController()
    sc := controller.NewSignupController()
    r.GET("/login", lc.LoginPage)
    r.POST("/login", lc.Login)
    r.GET("/logout", lc.Logout)
    r.GET("/signup", sc.SignupPage)
    r.POST("/signup", sc.Signup)


    ic := controller.NewIndexController()
    shc := controller.NewShiftController()
    wc := controller.NewWorkableController()
    pc := controller.NewProfileController()
    //ROOT (Authentication required)
    auth := r.Group("/")
    auth.Use(jwtauth.JwtAuthMiddleware())
    {
        auth.GET("/", ic.IndexPage)
        auth.GET("/shift/:year/:month", shc.ShiftPage)
        auth.POST("/shift/:year/:month", shc.Shift)
        auth.POST("/shift/:year/:month/auto-generate", shc.ShiftGenerate)
        auth.GET("/shift/:year/:month/storeholiday", shc.StoreHoliday)
        auth.GET("/workable/:year/:month", wc.WorkablePage)
        auth.POST("/workable/:year/:month", wc.Workable)
        auth.GET("/profile", pc.ProfilePage)
        auth.POST("/profile", pc.Profile)
    }

    //API
    api := r.Group("/api")
    {
        //API (Authentication required)
        apiAuth := api.Group("/")
        apiAuth.Use(jwtauth.JwtApiAuthMiddleware())
        {
            apiAuth.GET("/group/workables/:year/:month", controller.GroupWorkables)
            apiAuth.POST("/shift/:year/:month", controller.SaveShift)
        }
    }
    
    return r
}
