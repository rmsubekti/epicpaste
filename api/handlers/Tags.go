package handlers

import (
	"epicpaste/system/model"
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
