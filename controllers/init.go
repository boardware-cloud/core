package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boardware-cloud/common/server"
	"github.com/boardware-cloud/common/utils"
	api "github.com/boardware-cloud/core-api"
	"github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/middleware"
	model "github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
)

var router *gin.Engine

func Init() {
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	middleware.Health(router)
	api.AccountApiInterfaceMounter(router, accountApi)
	api.ServicesApiInterfaceMounter(router, serviceApi)
	api.VerificationApiInterfaceMounter(router, verificationApi)
	api.TicketApiInterfaceMounter(router, ticketApi)
	router.POST("/2account/webauthn/sessions/:id",
		func(ctx *gin.Context) {
			middleware.GetAccount(ctx,
				func(ctx *gin.Context, account model.Account) {
					var ccr protocol.CredentialCreationResponse
					if err := json.NewDecoder(ctx.Request.Body).Decode(&ccr); err != nil {
						return
					}
					fmt.Println(ccr.AttestationResponse)
					// id := ctx.Param("id")
					// if err := services.FinishRegistration(account, utils.StringToUint(id), "", "", ccr); err != nil {
					// 	err.GinHandler(ctx)
					// 	return
					// }
					// ctx.JSON(http.StatusCreated, "")
				})
		})
	router.POST("/2webauthn/sessions/tickets",
		func(ctx *gin.Context) {
			type request struct {
				Email string `json:"email"`
			}
			var req request
			ctx.ShouldBindJSON(&req)
			account, err := model.GetAccountByEmail(req.Email)
			if err != nil {
				err.GinHandler(ctx)
				return
			}
			option, session, err := services.BeginLogin(account)
			if err != nil {
				err.GinHandler(ctx)
				return
			}
			ctx.JSON(http.StatusCreated, gin.H{
				"id":        utils.UintToString(session.ID),
				"publicKey": option.Response,
			})
		})
	router.POST("/2webauthn/sessions/tickets/:id",
		func(ctx *gin.Context) {
			response, err := protocol.ParseCredentialRequestResponseBody(ctx.Request.Body)
			if err != nil {
				return
			}
			id := ctx.Param("id")
			ticket, errg := services.FinishLogin(utils.StringToUint(id), response)
			if errg != nil {
				errg.GinHandler(ctx)
				return
			}
			ctx.JSON(http.StatusCreated, api.Ticket{
				Token: ticket,
				Type:  api.WEBAUTHN,
			})
		})
}

func Run(addr ...string) {
	router.Run(addr...)
}
