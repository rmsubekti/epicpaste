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
	store := cookie.NewStore([]byte(os.Getenv("EPIC_COOKIE_SECRET_KEY")))
	v1.Use(sessions.Sessions("i92y", store))
	v1.Use(middleware.Auth())
	{
		v1.POST("/login", handler.UserLogin)
		v1.POST("/register", handler.UserRegister)
	}

	user := v1.Group("/:username")
	{
		user.GET("", handler.UserProfile)
		user.GET("/paste", handler.UserPastes)
	}

	paste := v1.Group("/paste")
	{
		paste.GET("", handler.ListPublicPaste)
		paste.POST("", handler.CreatePaste)
		paste.PATCH("/:id", handler.EditPaste)
		paste.GET("/:id", handler.ViewPaste)
		paste.DELETE("/:id", handler.DeletePaste)
	}

}
