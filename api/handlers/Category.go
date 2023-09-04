package handlers

import (
	"epicpaste/system/model"
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
