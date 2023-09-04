package handlers

import (
	"epicpaste/system/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code" swaggerignore:"true"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type LoginResponse struct {
	ID         string     `json:"id"`
	User       model.User `json:"user"`
	Token      string     `json:"token"`
	ExpireDays int        `json:"expire_days"`
	SignedDate time.Time  `json:"signed_date"`
}

func (r *Response) Json(c *gin.Context) {
	if r.Status == "" {
		r.Status = http.StatusText(r.Code)
	}
	c.JSON(r.Code, r)
}
