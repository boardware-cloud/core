package controllers

import (
	api "github.com/boardware-cloud/core-api"
	coreServices "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/middleware"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

var accountService = coreServices.GetAccountService()
var ticketService = coreServices.GetTicketService()
var verificationCodeService = coreServices.GetVerificationCodeService()

func init() {
	router = gin.Default()
	router.Use(accountService.Auth())
	router.Use(middleware.CorsMiddleware())
	middleware.Health(router)
	var accountApi AccountApi
	api.AccountApiInterfaceMounter(router, accountApi)
	api.VerificationApiInterfaceMounter(router, verificationApi)
	api.TicketApiInterfaceMounter(router, ticketApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
