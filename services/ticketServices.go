package services

import (
	"fmt"
	"strings"
	"time"

	errorCode "github.com/boardware-cloud/common/code"
	constants "github.com/boardware-cloud/common/constants/account"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model/core"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var ticketService *TicketService

func GetTicketService() *TicketService {
	if ticketService == nil {
		ticketService = NewTicketService()
	}
	return ticketService
}

func NewTicketService() *TicketService {
	return &TicketService{ticketRepository: core.GetTicketRepository()}
}

type TicketService struct {
	ticketRepository *core.TicketRepository
}

func (t TicketService) CreateTicket(email, ticketType string, password, verificationCode, totpCode *string) (string, error) {
	var account core.Account
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
			return t.createPasswordTicket(account, *password)
		}
	case "TOTP":
		if totpCode != nil {
			return t.createTotpTicket(account, *totpCode)
		}
	case "EMAIL":
		if verificationCode != nil {
			return t.createEmailTicket(account, *verificationCode)
		}
	}
	return "", errorCode.ErrUnauthorized
}

func TicketString(ticket core.Ticket) string {
	return fmt.Sprintf("%d:%s", ticket.ID, ticket.Secret)
}

func (t TicketService) UseTicket(token string) (core.Ticket, error) {
	ss := strings.Split(token, ":")
	if len(ss) != 2 {
		return core.Ticket{}, errorCode.ErrUnauthorized
	}
	ticket := t.ticketRepository.Find(utils.StringToUint(ss[0]))
	if ticket == nil {
		return core.Ticket{}, errorCode.ErrUnauthorized
	}
	if ticket.Secret != ss[1] {
		return core.Ticket{}, errorCode.ErrUnauthorized
	}
	DB.Delete(&ticket)
	if time.Now().Unix()-ticket.CreatedAt.Unix() > 300 {
		return core.Ticket{}, errorCode.ErrUnauthorized
	}
	return *ticket, nil
}

func (t TicketService) createTotpTicket(account core.Account, totpCode string) (string, error) {
	key, _ := otp.NewKeyFromURL(*account.Totp)
	if account.Totp == nil {
		return "", errorCode.ErrUnauthorized
	}
	if !totp.Validate(totpCode, key.Secret()) {
		return "", errorCode.ErrUnauthorized
	}
	return TicketString(ticketRepository.CreateTicket("TOTP", account.ID())), nil
}

func (t TicketService) createEmailTicket(account core.Account, verificationCode string) (string, error) {
	code := verificationCodeRepository.Get(account.Email, constants.TICKET)
	if !verify(code, verificationCode) {
		return "", errorCode.ErrVerificationCode
	}
	verificationCodeRepository.Delete(account.Email, constants.TICKET)
	return TicketString(ticketRepository.CreateTicket("EMAIL", account.ID())), nil
}

func (t TicketService) createPasswordTicket(account core.Account, password string) (string, error) {
	if !utils.PasswordsMatch(account.Password, password, account.Salt) {
		return "", errorCode.ErrUnauthorized
	}
	return TicketString(ticketRepository.CreateTicket("PASSWORD", account.ID())), nil
}
