package handlers

import (
	"epicpaste/system/model"
	u "epicpaste/system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePaste godoc
// @Summary Create a new paste
// @Description Currently login user can create a new paste
// @Tags         paste
// @Accept  json
// @Produce  json
// @Param request body model.Paste true "Body payload"
// @Success 200 {object} Response{data=model.Paste}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /paste [post]
func CreatePaste(c *gin.Context) {
	var (
		paste    model.Paste
		user     any
		ok       bool
		response Response
	)

	if user, ok = c.Get("user"); !ok || user == nil {
		response.Code = http.StatusUnauthorized
		response.Json(c)
		return
	}

	if err := c.ShouldBindJSON(&paste); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}
	paste.Paster = user.(model.User)

	if err := paste.Create(); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	response = Response{Code: http.StatusOK, Data: paste}
	response.Json(c)
}

// EditPaste godoc
// @Summary Edit a paste
// @Description Only owner can edit the paste
// @Tags         paste
// @Accept  json
// @Produce  json
// @Param        id   path      string  true  "Paste ID"
// @Success 200 {object} Response{data=model.Paste}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /paste/{id} [post]
func EditPaste(c *gin.Context) {
	var (
		paste    model.Paste
		user     any
		ok       bool
		response Response
	)

	if user, ok = c.Get("user"); !ok || user == nil {
		response.Code = http.StatusUnauthorized
		response.Json(c)
		return
	}

	if err := c.ShouldBindJSON(&paste); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	paste.ID = c.Param("id")
	paste.CreatedBy = user.(model.User).ID
	if err := paste.Update(); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	response = Response{Code: http.StatusOK, Data: paste}
	response.Json(c)
}

// DeletePaste godoc
// @Summary Delete a paste
// @Description Only owner can delete the paste
// @Tags         paste
// @Produce  json
// @Param        id   path      string  true  "Paste ID"
// @Success 200 {object} Response{data=model.Paste}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /paste/{id} [delete]
func DeletePaste(c *gin.Context) {
	var (
		paste    model.Paste
		user     any
		ok       bool
		response Response
	)

	if user, ok = c.Get("user"); !ok || user == nil {
		response.Code = http.StatusUnauthorized
		response.Json(c)
		return
	}

	paste.ID = c.Param("id")
	paste.CreatedBy = user.(model.User).ID
	if err := paste.Delete(); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	response = Response{
		Code:      http.StatusOK,
		Data:      paste,
		DeletedId: paste.ID,
	}
	response.Json(c)
}

// ViewPaste godoc
// @Summary View a paste
// @Description Paste can be viewed depending on visibility status of the paste
// @Tags         paste
// @Produce  json
// @Param        id   path      string  true  "Paste ID"
// @Success 200 {object} Response{data=model.Paste}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /paste/{id} [get]
func ViewPaste(c *gin.Context) {
	var paste model.Paste
	var response Response
	user, _ := c.Get("user")

	if err := paste.Get(c.Param("id")); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if !*paste.Public || (user != nil && paste.CreatedBy != user.(model.User).ID) {
		response.Code = http.StatusNotFound
		response.Json(c)
		return
	}

	response = Response{Code: http.StatusOK, Data: paste}
	response.Json(c)
}

// ListPublicPaste godoc
// @Summary View  list of pastes
// @Description Pastes can be viewed depending on visibility status of the paste
// @Tags         paste
// @Produce  json
// @Success 200 {object} Response{data=u.Paginator{rows=model.Pastes}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /paste [get]
func ListPublicPaste(c *gin.Context) {
	var paginator u.Paginator
	var pastes model.Pastes
	var response Response

	if err := c.Bind(&paginator); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if err := pastes.List(&paginator); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	response = Response{Code: http.StatusOK, Data: paginator}
	response.Json(c)
}

// UserPastes godoc
// @Summary View  list of pastes
// @Description Pastes can be viewed depending on logged in user
// @Tags         paste
// @Produce  json
// @Param        userId   path      string  true  "User ID"
// @Success 200 {object} Response{data=u.Paginator{rows=model.Pastes}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /{userId}/paste [get]
func UserPastes(c *gin.Context) {
	var paginator u.Paginator
	var pastes model.Pastes
	var response Response
	visitor, _ := c.Get("user")
	ownerId := c.Param("userId")

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
