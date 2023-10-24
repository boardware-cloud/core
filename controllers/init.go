package controllers

import (
	"net/http"

	api "github.com/boardware-cloud/core-api"
	coreServices "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine

var accountService coreServices.AccountService

func Init(inject *gorm.DB) {
	coreServices.Init(inject)
	accountService = coreServices.NewAccountService(inject)
	router = gin.Default()
	router.Use(accountService.Auth())
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://console-uat.k8s19932be1.boardware.com")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	middleware.Health(router)
	var accountApi AccountApi
	api.AccountApiInterfaceMounter(router, accountApi)
	api.VerificationApiInterfaceMounter(router, verificationApi)
	api.TicketApiInterfaceMounter(router, ticketApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
