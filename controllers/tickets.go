package controllers

import (
	"net/http"

	errorCode "github.com/boardware-cloud/common/code"
	coreapi "github.com/boardware-cloud/core-api"
	core "github.com/boardware-cloud/core/services"
	"github.com/gin-gonic/gin"
)

type TicketApi struct{}

var ticketApi TicketApi

// CreateTicket implements coreapi.TicketApiInterface.
func (TicketApi) CreateTicket(c *gin.Context, request coreapi.CreateTicketRequest) {
	ticketType := ""
	switch request.Type {
	case coreapi.WEBAUTHN, coreapi.EMAIL, coreapi.TOTP, coreapi.PASSWORD:
		ticketType = string(request.Type)
	default:
		c.JSON(http.StatusBadRequest, "")
		return
	}
	token, err := core.CreateTicket(request.Email, ticketType, request.Password, request.VerificationCode, request.TotpCode)
	if err != nil {
		errorCode.GinHandler(c, err)
		return
	}
	c.JSON(http.StatusCreated, coreapi.Ticket{
		Token: token,
		Type:  request.Type,
	})
}
