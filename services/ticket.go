package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/errors"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model/core"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func CreateTicket(email, ticketType string, password, verificationCode, totpCode *string) (string, *errors.Error) {
	var account core.Account
	ctx := DB.Where("email = ?", email).Find(&account)
	if ctx.RowsAffected == 0 {
		return "", errors.NotFoundError()
	}
	var loginRecord core.LoginRecord
	DB.Where("account_id = ?", account.ID).Order("created_at DESC").Limit(1).Find(&loginRecord)
	if time.Now().Unix()-loginRecord.CreatedAt.Unix() <= 1 {
		return "", errors.TooManyRequestsError()
	}
	DB.Save(&core.LoginRecord{AccountId: account.ID})
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
	return "", errors.AuthenticationError()
}

func createTotpTicket(account core.Account, totpCode string) (string, *errors.Error) {
	key, _ := otp.NewKeyFromURL(*account.Totp)
	if account.Totp == nil {
		return "", errors.AuthenticationError()
	}
	if !totp.Validate(totpCode, key.Secret()) {
		return "", errors.AuthenticationError()
	}
	return TicketString(createTicket("TOTP", account.ID)), nil
}

func createEmailTicket(account core.Account, verificationCode string) (string, *errors.Error) {
	code := GetVerification(account.Email, constants.TICKET)
	if !verify(code, verificationCode) {
		return "", errors.VerificationCodeError()
	}
	DB.Delete(&code)
	return TicketString(createTicket("EMAIL", account.ID)), nil
}

func createPasswordTicket(account core.Account, password string) (string, *errors.Error) {
	if !utils.PasswordsMatch(account.Password, password, account.Salt) {
		return "", errors.AuthenticationError()
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

func UseTicket(token string) (core.Ticket, *errors.Error) {
	ss := strings.Split(token, ":")
	var ticket core.Ticket
	if len(ss) != 2 {
		return ticket, errors.AuthenticationError()
	}
	ctx := DB.Find(&ticket, utils.StringToUint(ss[0]))
	if ctx.RowsAffected == 0 {
		return ticket, errors.AuthenticationError()
	}
	if ticket.Secret != ss[1] {
		return ticket, errors.AuthenticationError()
	}
	DB.Delete(&ticket)
	if time.Now().Unix()-ticket.CreatedAt.Unix() > 300 {
		return ticket, errors.AuthenticationError()
	}
	return ticket, nil
}
