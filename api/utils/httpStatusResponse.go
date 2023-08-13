package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONErr(code int, c *gin.Context, message any) {
	var response = make(map[string]any)
	response["code"] = code
	response["status"] = http.StatusText(code)
	if message != nil {
		response["message"] = message
	}
	c.JSON(code, response)
}
