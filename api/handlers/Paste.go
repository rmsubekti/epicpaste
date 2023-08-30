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
// @Security Bearer
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
// @Param request body model.Paste true "Body payload"
// @Param        id   path      string  true  "Paste ID"
// @Success 200 {object} Response{data=model.Paste}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /paste/{id} [patch]
// @Security Bearer
func EditPaste(c *gin.Context) {
	var (
		paste    model.Paste
		user     any
		ok       bool
		response Response
	)

	if user, ok = c.Get("user"); !ok || user == nil {
		response = Response{
			Code:    http.StatusUnauthorized,
			Message: "Please login first",
		}
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
// @Security Bearer
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
		Code:    http.StatusOK,
		Data:    paste,
		Message: "Deleted successfully",
	}
	response.Json(c)
}

// ViewPaste godoc
// @Summary View a paste
// @Description Paste can be viewed depending on visibility status of the paste.
// @Description Bearer Token is Optional
// @Tags         paste
// @Produce  json
// @Param        id   path      string  true  "Paste ID"
// @Success 200 {object} Response{data=model.Paste}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /paste/{id} [get]
// @Security Bearer
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
// @Description Bearer Token is Optional
// @Tags         paste
// @Produce  json
// @Param        page    query    int  false  "show data on page n"
// @Param        limit    query     int  false  "limit items per page"
// @Param        q    query     string  false  "filter query"
// @Success 200 {object} Response{data=u.Paginator{rows=model.Pastes}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /paste [get]
// @Security Bearer
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
