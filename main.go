package main

import (
	"epicpaste/api"
	_ "epicpaste/docs"
	"epicpaste/proto"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Epic Paste Service
// @version         1.0
// @description     A snippet management service API in Go using Gin framework.
// @termsOfService  https://bekti.net

// @contact.name   Rahmat Subekti
// @contact.url    https://twitter.com/rmsiannnaksnfe
// @contact.email  rmssssssgmail@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3030
// @BasePath  /v1
func main() {
	app := gin.Default()
	PORT := os.Getenv("PORT")
	GRPC := os.Getenv("GRPC")

	if GRPC == "true" {
		go proto.Start()
	}

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	app.Use(cors.New(config))

	api.Serve(app)

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Message": "Status OK",
		})
	})
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Run(":" + PORT)

}
