package services

import (
	errorCode "github.com/boardware-cloud/common/code"
	"github.com/boardware-cloud/model/core"
	"gorm.io/gorm"
)

func NewTicketService(db *gorm.DB) TicketService {
	return TicketService{ticketRepository: core.NewTicketRepository(db)}
}

type TicketService struct {
	ticketRepository core.TicketRepository
}

func (t TicketService) CreateTicket(email, ticketType string, password, verificationCode, totpCode *string) (string, error) {
	var account core.Account
	// t.ticketRepository.Find()
	ctx := DB.Where("email = ?", email).Find(&account)
	if ctx.RowsAffected == 0 {
		return "", errorCode.ErrNotFound
	}
	if !account.ColdDown(500) {
		return "", errorCode.ErrTooManyRequests
	}
	account.CreateColdDown()
	switch ticketType {
	case "PASSWORD":
		if password != nil {
			return createPasswordTicket(account, *password)
		}
	case "TOTP":
		if totpCode != nil {
			return createTotpTicket(account, *totpCode)
		}
	case "EMAIL":
		if verificationCode != nil {
			return createEmailTicket(account, *verificationCode)
		}
	}
	return "", errorCode.ErrUnauthorized
}
