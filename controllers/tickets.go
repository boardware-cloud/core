package controllers

import (
	"net/http"

	coreapi "github.com/boardware-cloud/core-api"
	core "github.com/boardware-cloud/core/services"
	"github.com/chenyunda218/golambda"
	"github.com/gin-gonic/gin"
)

type TicketApi struct{}

var ticketApi TicketApi

// CreateTicket implements coreapi.TicketApiInterface.
func (TicketApi) CreateTicket(c *gin.Context, request coreapi.CreateTicketRequest) {
	token, err := core.CreateTicket(request.Email, string(request.Type), request.Password, request.VerificationCode, request.TotpCode)
	if err != nil {
		err.GinHandler(c)
		return
	}
	c.JSON(http.StatusCreated, coreapi.Ticket{
		Token: golambda.Reference(token),
		Type:  golambda.Reference(request.Type),
	})
}
