package api

import (
	handler "epicpaste/api/handlers"
	"epicpaste/api/middleware"
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
		auth.POST("/login", handler.UserLogin)
		auth.POST("/register", handler.UserRegister)
	}

	user := v1.Group("/:userId")
	{
		user.GET("", handler.UserPastes)
		user.GET("/paste", handler.UserPastes)
	}

	paste := v1.Group("/paste")
	{
		paste.GET("", handler.ListPublicPaste)
		paste.POST("", handler.CreatePaste)
		paste.POST("/:id", handler.EditPaste)
		paste.GET("/:id", handler.ViewPaste)
		paste.DELETE("/:id", handler.DeletePaste)
	}

}
