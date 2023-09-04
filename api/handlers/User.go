package handlers

import (
	"epicpaste/system/model"
	u "epicpaste/system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserProfile godoc
// @Summary View  user profile
// @Description View user profile
// @Description Bearer Token is Optional, Use bearer token to see private user
// @Tags         user
// @Produce  json
// @Param        username   path  string  true  "UserName" example(epicpaster)
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
// @Summary View  list of paste
// @Description Pastes can be viewed depending on logged in user
// @Description Bearer Token is Optional, Use bearer token to see private user pastes
// @Tags         user
// @Produce  json
// @Param        username   path      string  true  "UserName" example(epicpaster)
// @Success 200 {object} Response{data=u.Paginator{rows=model.Pastes}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /{username}/paste [get]
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

	if visitor != nil && visitor.(model.User).UserName == ownerId {
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

// EditProfile godoc
// @Summary Edit User Profile
// @Description Username cannot be changed
// @Tags         user
// @Produce  json
// @Param        username   path      string  true  "UserName" example(epicpaster)
// @Param request body string true " Body payload message/rfc822" SchemaExample({\n\t"name": "Epic Paster"\n})
// @Success 200 {object} Response{data=model.User}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /{username} [patch]
// @Security Bearer
func EditProfile(c *gin.Context) {
	var user model.User
	var editor any
	var response Response
	var ok bool

	if editor, ok = c.Get("user"); !ok || editor == nil {
		response = Response{
			Code:    http.StatusUnauthorized,
			Message: "Please login first",
		}
		response.Json(c)
		return
	}

	if editor.(model.User).UserName != c.Param("username") {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: "You can not edit other user profile",
		}
		response.Json(c)
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	user.UserName = c.Param("username")

	if err := user.Update(); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}
	response = Response{Code: http.StatusOK, Data: user}
	response.Json(c)
}
