package controllers

import (
	api "github.com/boardware-cloud/core-api"
	"github.com/boardware-cloud/middleware"
	model "github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine
var db *gorm.DB
var accountRepository model.AccountRepository

func Init(inject *gorm.DB) {
	db = inject
	model.Init(db)
	accountRepository = model.NewAccountRepository(db)
	router = gin.Default()
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
