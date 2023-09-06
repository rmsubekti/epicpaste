package handlers

import (
	"epicpaste/system/model"
	u "epicpaste/system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListCategory godoc
// @Summary View  all categories
// @Tags         taxonomy
// @Produce  json
// @Success 200 {object} Response{data=model.Categories}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /category [get]
func ListCategory(c *gin.Context) {
	var categories model.Categories
	var response Response
	if err := categories.List(); err != nil {
		response = Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}
	response = Response{Code: http.StatusOK, Data: categories}
	response.Json(c)
}

// ListPasteByCategory godoc
// @Summary Paste By Category
// @Description Pastes can be viewed depending on visibility status of the paste
// @Tags         taxonomy
// @Produce  json
// @Param        page    query    int  false  "show data on page n"
// @Param        limit    query     int  false  "limit items per page"
// @Param        q    query     string  false  "filter query"
// @Param        category   path      string  true  "Category name"
// @Success 200 {object} Response{data=u.Paginator{rows=model.Pastes}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router /category/{category} [get]
func ListPasteByCategory(c *gin.Context) {
	var response Response
	var pastes model.Pastes
	var paginator u.Paginator
	category := c.Param("category")

	if err := c.Bind(&paginator); err != nil {
		response = Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		response.Json(c)
		return
	}

	if err := pastes.ListByCategory(category, &paginator); err != nil {
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
