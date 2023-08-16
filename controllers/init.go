package controllers

import (
	"github.com/boardware-cloud/common/server"
	api "github.com/boardware-cloud/core-api"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() {
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	api.AccountApiInterfaceMounter(router, accountApi)
	api.ServicesApiInterfaceMounter(router, serviceApi)
	api.VerificationApiInterfaceMounter(router, verificationApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
