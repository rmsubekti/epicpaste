package handlers

import (
	"epicpaste/system/auth"
	"epicpaste/system/model"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AccountLogin godoc
// @Summary Login
// @Description Create new session
// @Description user can fill the username field with their username or email
// @Tags         account
// @Produce  json
// @Param request body string true " Body payload message/rfc822" SchemaExample({\n\t"username": "epicpaster",\n\t"password": "5uperSecret"\n})
// @Success 200 {object} Response{data=LoginResponse}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /login [post]
func AccountLogin(c *gin.Context) {
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
	session.Options(sessions.Options{MaxAge: (24 * 60 * 60 * sessionDays)})
	session.Save()

	response = Response{
		Code: http.StatusOK,
		Data: LoginResponse{
			ID:         account.ID,
			User:       account.User,
			Token:      token,
			ExpireDays: sessionDays,
			SignedDate: time.Now(),
		},
	}
	response.Json(c)
}

// AccountRegister godoc
// @Summary Register a new account
// @Description Usename cannot be changed after account is created
// @Tags         account
// @Produce  json
// @Param request body model.Account true "Body payload"
// @Success 200 {object} Response{data=bool}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /register [post]
func AccountRegister(c *gin.Context) {
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

// PasswordChange godoc
// @Summary Change Account password
// @Tags         account
// @Produce  json
// @Param request body model.ChangePassword true "Body payload"
// @Description Need to login first
// @Success 200 {object} Response
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /change-password [patch]
// @Security Bearer
func PasswordChange(c *gin.Context) {
	var (
		user     any
		ok       bool
		response Response
		password model.ChangePassword
	)

	if user, ok = c.Get("user"); !ok || user == nil {
		response = Response{
			Code:    http.StatusUnauthorized,
			Message: "Please login first",
		}
		response.Json(c)
		return
	}

	if err := c.ShouldBindJSON(&password); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	account := &model.Account{UserName: user.(model.User).UserName}
	if err := account.ChangePassword(password); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	response = Response{
		Code:    http.StatusOK,
		Message: "Password changed successfuly",
	}
	response.Json(c)
}

// ChangeEmail godoc
// @Summary Change Account email
// @Tags         account
// @Produce  json
// @Param request body Emailpayload true "Body payload"
// @Description Need to login first
// @Success 200 {object} Response
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /change-email [patch]
// @Security Bearer
func ChangeEmail(c *gin.Context) {
	var (
		user     any
		ok       bool
		response Response
		account  model.Account
		email    Emailpayload
	)

	if user, ok = c.Get("user"); !ok || user == nil {
		response = Response{
			Code:    http.StatusUnauthorized,
			Message: "Please login first",
		}
		response.Json(c)
		return
	}

	if err := c.ShouldBindJSON(&email); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if err := account.Get(user.(model.User).UserName); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}
	if err := account.ChangeEmail(email.Email); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	response = Response{
		Code:    http.StatusOK,
		Message: "email changed",
		Data:    email,
	}
	response.Json(c)
}

// LogOut godoc
// @Summary User Account Logout
// @Description This only work with cookie.
// @Description For JWT Token, you must set token from the respose to the frontend.
// @Description Need to login first
// @Tags         account
// @Produce  json
// @Success 200 {object} Response
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /logout [get]
// @Security Bearer
func LogOut(c *gin.Context) {
	var (
		user        any
		ok          bool
		response    Response
		sessionDays int = -1
	)

	if user, ok = c.Get("user"); !ok || user == nil {
		response = Response{
			Code:    http.StatusUnauthorized,
			Message: "Please login first",
		}
		response.Json(c)
		return
	}

	token, err := auth.CreateLoginSignature(user.(*model.User), sessionDays)
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
	session.Options(sessions.Options{MaxAge: (24 * 60 * 60 * sessionDays)})
	session.Save()
	response = Response{
		Code:    http.StatusOK,
		Message: "Logged out",
		Data:    token,
	}
	response.Json(c)
}
