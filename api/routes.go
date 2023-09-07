package api

import (
	handler "epicpaste/api/handlers"
	"epicpaste/api/middleware"
	"epicpaste/system/helper"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Serve(app *gin.Engine) {
	v1 := app.Group("/v1")
	//cors setting
	store := cookie.NewStore([]byte(helper.GetEnv("EPIC_COOKIE_SECRET_KEY", "epic_cookie")))
	v1.Use(sessions.Sessions("i92y", store))
	v1.Use(middleware.Auth())
	{
		v1.POST("/login", handler.AccountLogin)
		v1.POST("/register", handler.AccountRegister)
		v1.PATCH("/change-password", handler.PasswordChange)
		v1.PATCH("/change-email", handler.ChangeEmail)
		v1.GET("/logout", handler.LogOut)
	}

	user := v1.Group("/:username")
	{
		user.GET("", handler.UserProfile)
		user.PATCH("", handler.EditProfile)
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

	category := v1.Group("/category")
	{
		category.GET("", handler.ListCategory)
		category.GET("/:category", handler.ListPasteByCategory)
	}

	syntax := v1.Group("/syntax")
	{
		syntax.GET("", handler.ListSyntax)
		syntax.GET("/:syntax", handler.ListPasteBySyntax)
	}

	tag := v1.Group("/tag")
	{
		tag.GET("", handler.ListTags)
		tag.GET("/:tag", handler.ListPasteByTag)
	}

}
