package main

import (
	"epicpaste/api"
	"epicpaste/docs"
	"epicpaste/proto"
	"epicpaste/system/helper"
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
// @description     Login to create token.

// @contact.name   Rahmat Subekti
// @contact.url    https://bekti.net/social
// @contact.email  rmsubekti2011@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath  /v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	gin.SetMode(helper.GetEnv("EPIC_DEBUG", gin.ReleaseMode))
	app := gin.Default()
	PORT := helper.GetEnv("EPIC_PORT", "80")
	HOSTNAME := helper.GetEnv("EPIC_HOSTNAME", "localhost")

	docs.SwaggerInfo.Host = HOSTNAME + ":" + PORT

	if os.Getenv("EPIC_GRPC") == "true" {
		go proto.Start()
	}

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PATCH", "GET", "DELETE"},
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

	app.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		//Convert your cached html string to byte array
		c.Writer.Write([]byte(`<a style="font-size:2em" href="/docs/index.html">Open the documentation</a>`))
	})
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Run(":" + PORT)

}
