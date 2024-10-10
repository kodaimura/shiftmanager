package server

import (
	"github.com/gin-gonic/gin"

	"shiftmanager/internal/middleware"
	"shiftmanager/internal/controller"
)

/*
 Routing for "/" 
*/
func SetWebRouter(r *gin.RouterGroup) {
	ic := controller.NewIndexController()
	ac := controller.NewAccountController()
	apc := controller.NewAccountProfileController()
	spc := controller.NewShiftPreferredController()

	r.GET("/signup", middleware.BasicAuthMiddleware(), ac.SignupPage)
	r.GET("/login", ac.LoginPage)
	r.GET("/logout", ac.Logout)

	auth := r.Group("", middleware.JwtAuthMiddleware())
	{
		auth.GET("/", ic.IndexPage)
		auth.GET("/account_profile", apc.AccountProfilePage)
		auth.GET("/shift_preferred/:year/:month", spc.ShiftPreferredPage)
	}
}


/*
 Routing for "/api"
*/
func SetApiRouter(r *gin.RouterGroup) {
	ac := controller.NewAccountController()
	apc := controller.NewAccountProfileController()
	spc := controller.NewShiftPreferredController()

	r.POST("/signup", ac.ApiSignup)
	r.POST("/login", ac.ApiLogin)

	auth := r.Group("", middleware.JwtAuthApiMiddleware())
	{
		auth.GET("/account", ac.ApiGetOne)
		auth.PUT("/account/name", ac.ApiPutName)
		auth.PUT("/account/password", ac.ApiPutPassword)
		auth.DELETE("/account", ac.ApiDelete)

		auth.GET("/account_profile", apc.ApiGetOne)
		auth.POST("/account_profile", apc.ApiPost)

		auth.GET("/shift_preferred", spc.ApiGet)
		auth.GET("/shift_preferred/:year/:month", spc.ApiGetOne)
		auth.POST("/shift_preferred/:year/:month", spc.ApiPost)
	}
}