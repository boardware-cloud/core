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

func CreateTicket(email, ticketType string, password, verificationCode, totpCode *string) (string, error) {
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

func createTotpTicket(account core.Account, totpCode string) (string, error) {
	key, _ := otp.NewKeyFromURL(*account.Totp)
	if account.Totp == nil {
		return "", errorCode.ErrUnauthorized
	}
	if !totp.Validate(totpCode, key.Secret()) {
		return "", errorCode.ErrUnauthorized
	}
	return TicketString(createTicket("TOTP", account.ID)), nil
}

func createEmailTicket(account core.Account, verificationCode string) (string, error) {
	code := GetVerification(account.Email, constants.TICKET)
	if !verify(code, verificationCode) {
		return "", errorCode.ErrVerificationCode
	}
	DB.Delete(&code)
	return TicketString(createTicket("EMAIL", account.ID)), nil
}

func createPasswordTicket(account core.Account, password string) (string, error) {
	if !utils.PasswordsMatch(account.Password, password, account.Salt) {
		return "", errorCode.ErrUnauthorized
	}
	return TicketString(createTicket("PASSWORD", account.ID)), nil
}

func createTicket(t string, accountId uint) core.Ticket {
	var ticket core.Ticket
	ticket.AccountId = accountId
	ticket.Secret = RandomNumberString(16)
	ticket.Type = t
	DB.Save(&ticket)
	return ticket
}

func TicketString(ticket core.Ticket) string {
	return fmt.Sprintf("%d:%s", ticket.ID, ticket.Secret)
}

func UseTicket(token string) (core.Ticket, error) {
	ss := strings.Split(token, ":")
	var ticket core.Ticket
	if len(ss) != 2 {
		return ticket, errorCode.ErrUnauthorized
	}
	ctx := DB.Find(&ticket, utils.StringToUint(ss[0]))
	if ctx.RowsAffected == 0 {
		return ticket, errorCode.ErrUnauthorized
	}
	if ticket.Secret != ss[1] {
		return ticket, errorCode.ErrUnauthorized
	}
	DB.Delete(&ticket)
	if time.Now().Unix()-ticket.CreatedAt.Unix() > 300 {
		return ticket, errorCode.ErrUnauthorized
	}
	return ticket, nil
}
