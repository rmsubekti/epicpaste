package handlers

import (
	"epicpaste/api/utils"
	"epicpaste/system/auth"
	"epicpaste/system/model"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
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

func UserRegister(c *gin.Context) {
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

func UserProfile(c *gin.Context) {
	var account model.Account
	id := c.Param("id")
	visitor, _ := c.Get("user")

	if err := account.Get(id); err != nil {
		utils.JSONErr(http.StatusNotFound, c, nil)
		return
	}

	if !account.Setting.Crawlable && visitor == nil {
		utils.JSONErr(http.StatusNotFound, c, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": account.User})
}
