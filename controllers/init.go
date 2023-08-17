package controllers

import (
	"net/http"

	"github.com/boardware-cloud/common/server"
	api "github.com/boardware-cloud/core-api"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

type health struct {
	Status string   `json:"status"`
	Checks []string `json:"checks"`
}

func Init() {
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	router.GET("/health/ready", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, health{
			Status: "UP",
			Checks: make([]string, 0),
		})
	})
	router.GET("/health/live", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, health{
			Status: "UP",
			Checks: make([]string, 0),
		})
	})
	api.AccountApiInterfaceMounter(router, accountApi)
	api.ServicesApiInterfaceMounter(router, serviceApi)
	api.VerificationApiInterfaceMounter(router, verificationApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
