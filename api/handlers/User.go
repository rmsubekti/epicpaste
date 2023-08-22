package handlers

import (
	"epicpaste/system/auth"
	"epicpaste/system/model"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	var account model.Account
	var response Response
	const sessionDays = 1
	if err := c.ShouldBindJSON(&account); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	//verify login
	if err := account.Login(); err != nil {
		response = Response{
			Code:    http.StatusNotAcceptable,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	token, err := auth.CreateLoginSignature(&account.User, sessionDays)
	if err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}
	session := sessions.Default(c)
	session.Set("token", "Bearer "+token)
	session.Save()

	user := map[string]any{
		"id":        account.ID,
		"user_name": account.UserName,
		"full_name": account.User.Name,
	}

	response = Response{
		Code: http.StatusOK,
		Data: gin.H{
			"user":        user,
			"token":       token,
			"expire_days": sessionDays,
			"signed_date": time.Now(),
		},
	}
	response.Json(c)
}

func UserRegister(c *gin.Context) {
	var account model.Account
	var response Response
	if err := c.ShouldBindJSON(&account); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if err := account.Register(); err != nil {
		response = Response{
			Code:    http.StatusNotAcceptable,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	response = Response{
		Code:    http.StatusOK,
		Message: "User registered",
		Data:    true,
	}
	response.Json(c)
}

func UserProfile(c *gin.Context) {
	var account model.Account
	var response Response
	id := c.Param("id")
	visitor, _ := c.Get("user")

	if err := account.Get(id); err != nil {
		response = Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if !account.Setting.Crawlable && visitor == nil {
		response = Response{
			Code: http.StatusNotFound,
		}
		response.Json(c)
		return
	}

	response = Response{
		Code: http.StatusOK,
		Data: account.User,
	}
	response.Json(c)
}
