package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/handlers"
	"github.com/muplat/muplat-backend/pkg/setup"
)

func main() {
	log.Println("Starting server...")
	globalConf := setup.InitGlobalConfig()
	httpHandler := handlers.NewHttpHandler(globalConf.Db, globalConf.ClusterConn, globalConf.Jwt)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.RedirectTrailingSlash = true

	DeclareRoutes(r, httpHandler)

	r.Run(":8080")
}
