package controllers

import (
	api "github.com/boardware-cloud/core-api"
	coreServices "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine

var accountService coreServices.AccountService
var ticketService coreServices.TicketService

var verificationCodeService coreServices.VerificationCodeService

func Init(inject *gorm.DB) {
	coreServices.Init(inject)
	accountService = coreServices.NewAccountService(inject)
	ticketService = coreServices.NewTicketService(inject)
	verificationCodeService = coreServices.NewVerificationCodeService(inject)
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
