package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/controllers"
	"github.com/muplat/muplat-backend/pkg/setup"
	"github.com/muplat/muplat-backend/services"
)

func main() {
	globalConfig := setup.InitGlobalConfig()
	HttpHandler := &controllers.HttpHandler{
		AppService: &services.AppService{},
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.RedirectTrailingSlash = true

	DeclareRoutes(r)

	r.Run(":8080")
}
