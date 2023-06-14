package app

import (
	"epicpaste/app/controller"
	"epicpaste/app/middleware"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Serve(app *gin.Engine) {
	v1 := app.Group("/v1")
	//cors setting
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	v1.Use(cors.New(config))

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

	note := v1.Group("/paste")
	{
		note.GET("", controller.ListPublicPaste)
		note.POST("", controller.CreatePaste)
		note.POST("/:id", controller.EditPaste)
		note.GET("/:id", controller.ViewPaste)
	}

}
