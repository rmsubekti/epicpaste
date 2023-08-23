package handlers

import (
	"epicpaste/system/auth"
	"epicpaste/system/model"
	u "epicpaste/system/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserLogin godoc
// @Summary User Login
// @Description Create new user session
// @Tags         user
// @Produce  json
// @Param request body string true " Body payload message/rfc822" SchemaExample({\n\t"username": "epicpaster",\n\t"password": "5uperSecret"\n})
// @Success 200 {object} Response{data=LoginResponse}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /login [post]
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

	response = Response{
		Code: http.StatusOK,
		Data: LoginResponse{
			ID:         account.ID,
			UserName:   account.UserName,
			Name:       account.User.Name,
			Token:      token,
			ExpireDays: sessionDays,
			SignedDate: time.Now(),
		},
	}
	response.Json(c)
}

// UserRegister godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags         user
// @Produce  json
// @Param request body model.Account true "Body payload"
// @Success 200 {object} Response{data=bool}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /register [post]
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

// UserProfile godoc
// @Summary View  user profile
// @Description View user profile
// @Description Bearer Token is Optional, Use bearer token to see private user
// @Tags         user
// @Produce  json
// @Param        userId   path      string  true  "User ID"
// @Success 200 {object} Response{data=model.User}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /{username} [get]
// @Security Bearer
func UserProfile(c *gin.Context) {
	var account model.Account
	var response Response
	username := c.Param("username")
	visitor, _ := c.Get("user")

	if err := account.Get(username); err != nil {
		response = Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if !account.Setting.Crawlable && visitor == nil {
		response = Response{
			Code: http.StatusForbidden,
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

// UserPastes godoc
// @Summary View  list of pastes
// @Description Pastes can be viewed depending on logged in user
// @Description Bearer Token is Optional, Use bearer token to see private user pastes
// @Tags         user
// @Produce  json
// @Param        userId   path      string  true  "User ID"
// @Success 200 {object} Response{data=u.Paginator{rows=model.Pastes}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /{userId}/paste [get]
// @Security Bearer
func UserPastes(c *gin.Context) {
	var paginator u.Paginator
	var pastes model.Pastes
	var response Response
	visitor, _ := c.Get("user")
	ownerId := c.Param("username")

	if err := c.Bind(&paginator); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if visitor != nil && visitor.(model.User).ID == ownerId {
		// all pastes by user
		if err := pastes.ListByUser(ownerId, false, &paginator); err != nil {
			response = Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			response.Json(c)
			return
		}
	} else {
		// list all public
		if err := pastes.ListByUser(ownerId, true, &paginator); err != nil {
			response = Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			response.Json(c)
			return
		}
	}

	response = Response{Code: http.StatusOK, Data: paginator}
	response.Json(c)
}
