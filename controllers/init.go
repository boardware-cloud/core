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

func Init(inject *gorm.DB) {
	coreServices.Init(inject)
	accountService = coreServices.NewAccountService(inject)
	ticketService = coreServices.NewTicketService(inject)
	router = gin.Default()
	router.Use(accountService.Auth())
	// router.Use(func(c *gin.Context) {
	// 	method := c.Request.Method
	// 	origin := c.Request.Header.Get("Origin")
	// 	if origin != "" {
	// 		c.Header("Access-Control-Allow-Origin", "*")
	// 		c.Header("Access-Control-Allow-Methods", "*")
	// 		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	// 		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	// 		c.Header("Access-Control-Allow-Credentials", "true")
	// 	}
	// 	if method == "OPTIONS" {
	// 		c.AbortWithStatus(http.StatusNoContent)
	// 	}
	// 	c.Next()
	// })
	// router.Use(middleware.CorsMiddleware())
	middleware.Health(router)
	var accountApi AccountApi
	api.AccountApiInterfaceMounter(router, accountApi)
	api.VerificationApiInterfaceMounter(router, verificationApi)
	api.TicketApiInterfaceMounter(router, ticketApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
