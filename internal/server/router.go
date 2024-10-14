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
		auth.GET("/account_profiles/me", apc.AccountProfilePage)
		auth.GET("/shift_preferreds/me/:year/:month", spc.ShiftPreferredPage)
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
		auth.GET("/accounts/me", ac.ApiGetOne)
		auth.PUT("/accounts/me/name", ac.ApiPutName)
		auth.PUT("/accounts/me/password", ac.ApiPutPassword)
		auth.DELETE("/accounts/me", ac.ApiDelete)

		auth.GET("/account_profiles", apc.ApiGet)
		auth.GET("/account_profiles/me", apc.ApiGetOne)
		auth.POST("/account_profiles/me", apc.ApiPost)

		auth.GET("/shift_preferreds", spc.ApiGet)
		auth.GET("/shift_preferreds/me/:year/:month", spc.ApiGetOne)
		auth.POST("/shift_preferreds/me/:year/:month", spc.ApiPost)
	}
}