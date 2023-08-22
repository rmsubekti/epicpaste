package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
	DeletedId any    `json:"deletedId,omitempty"`
}

func (r *Response) Json(c *gin.Context) {
	if r.Status == "" {
		r.Status = http.StatusText(r.Code)
	}
	c.JSON(r.Code, r)
}
