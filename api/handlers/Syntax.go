package handlers

import (
	"epicpaste/system/model"
	u "epicpaste/system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListSyntax godoc
// @Summary View  all syntaxs
// @Tags         taxonomy
// @Produce  json
// @Success 200 {object} Response{data=model.Syntaxs}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /syntax [get]
func ListSyntax(c *gin.Context) {
	var syntaxs model.Syntaxs
	var response Response
	if err := syntaxs.List(); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}
	response = Response{Code: http.StatusOK, Data: syntaxs}
	response.Json(c)
}

// ListPasteBySyntax godoc
// @Summary Paste By Syntax
// @Description Pastes can be viewed depending on visibility status of the paste
// @Tags         taxonomy
// @Produce  json
// @Param        page    query    int  false  "show data on page n"
// @Param        limit    query     int  false  "limit items per page"
// @Param        q    query     string  false  "filter query"
// @Param        syntax   path      string  true  "Syntax name"
// @Success 200 {object} Response{data=u.Paginator{rows=model.Pastes}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /syntax/{syntax} [get]
func ListPasteBySyntax(c *gin.Context) {
	var response Response
	var pastes model.Pastes
	var paginator u.Paginator
	syntax := c.Param("syntax")

	if err := c.Bind(&paginator); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if err := pastes.ListBySyntax(syntax, &paginator); err != nil {
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
