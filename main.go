package main

import (
	webapp "epicpaste/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	webapp.Serve(app)

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Message": "Status OK",
		})
	})
	app.Run()
}
