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
	rc := controller.NewRootController()
	ac := controller.NewAccountController()

	r.GET("/signup", ac.SignupPage)
	r.GET("/login", ac.LoginPage)
	r.GET("/logout", ac.Logout)

	auth := r.Group("", middleware.JwtAuthMiddleware())
	{
		auth.GET("/", rc.IndexPage)
	}
}


/*
 Routing for "/api"
*/
func SetApiRouter(r *gin.RouterGroup) {
	ac := controller.NewAccountController()

	r.POST("/signup", ac.ApiSignup)
	r.POST("/login", ac.ApiLogin)

	auth := r.Group("", middleware.JwtAuthApiMiddleware())
	{
		auth.GET("/account/profile", ac.ApiGetProfile)
		auth.PUT("/account/name", ac.ApiPutName)
		auth.PUT("/account/password", ac.ApiPutPassword)
		auth.DELETE("/account", ac.ApiDeleteAccount)
	}
}