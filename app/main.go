package app

import (
	"epicpaste/app/controller"
	"epicpaste/app/middleware"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Serve(app *gin.Engine) {
	v1 := app.Group("/v1")
	//cors setting
	store := cookie.NewStore([]byte(os.Getenv("API_SESSION_KEY")))
	v1.Use(sessions.Sessions("i92y", store))
	v1.Use(middleware.Auth())

	auth := v1.Group("/auth")
	{
		auth.POST("/login", controller.UserLogin)
		auth.POST("/register", controller.UserRegister)
	}

	user := v1.Group("/:userId")
	{
		user.GET("", controller.UserPastes)
		user.GET("/paste", controller.UserPastes)
	}

	paste := v1.Group("/paste")
	{
		paste.GET("", controller.ListPublicPaste)
		paste.POST("", controller.CreatePaste)
		paste.POST("/:id", controller.EditPaste)
		paste.GET("/:id", controller.ViewPaste)
		paste.DELETE("/:id", controller.DeletePaste)
	}

}
