package controller

import (
	"epicpaste/app/utils"
	"epicpaste/system/auth"
	"epicpaste/system/model"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var UserLogin = func(c *gin.Context) {
	var account model.Account
	const sessionDays = 1
	if err := c.ShouldBindJSON(&account); err != nil {
		utils.JSONErr(http.StatusBadRequest, c, err.Error())
		return
	}

	//verify login
	if err := account.Login(); err != nil {
		utils.JSONErr(http.StatusNotAcceptable, c, err.Error())
		return
	}

	token, err := auth.CreateLoginSignature(&account.User, sessionDays)
	if err != nil {
		utils.JSONErr(http.StatusInternalServerError, c, err.Error())
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

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token, "expire_days": sessionDays, "signed_date": time.Now()})
}

var UserRegister = func(c *gin.Context) {
	var account model.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		utils.JSONErr(http.StatusBadRequest, c, err.Error())
		return
	}

	if err := account.Register(); err != nil {
		utils.JSONErr(http.StatusNotAcceptable, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
