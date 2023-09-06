package handlers

import (
	"epicpaste/system/model"
	u "epicpaste/system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListTags godoc
// @Summary View  all tags
// @Tags         taxonomy
// @Produce  json
// @Success 200 {object} Response{data=model.Tags}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /tag [get]
func ListTags(c *gin.Context) {
	var tags model.Tags
	var response Response
	if err := tags.List(); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}
	response = Response{Code: http.StatusOK, Data: tags}
	response.Json(c)
}

// ListPasteByTag godoc
// @Summary Paste By Tag
// @Description Pastes can be viewed depending on visibility status of the paste
// @Tags         taxonomy
// @Produce  json
// @Param        page    query    int  false  "show data on page n"
// @Param        limit    query     int  false  "limit items per page"
// @Param        q    query     string  false  "filter query"
// @Param        tag   path      string  true  "Tag name"
// @Success 200 {object} Response{data=u.Paginator{rows=model.Pastes}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /tag/{tag} [get]
func ListPasteByTag(c *gin.Context) {
	var response Response
	var pastes model.Pastes
	var paginator u.Paginator
	tag := c.Param("tag")

	if err := c.Bind(&paginator); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if err := pastes.ListByTag(tag, &paginator); err != nil {
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
