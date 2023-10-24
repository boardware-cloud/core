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

func createTotpTicket(account core.Account, totpCode string) (string, error) {
	key, _ := otp.NewKeyFromURL(*account.Totp)
	if account.Totp == nil {
		return "", errorCode.ErrUnauthorized
	}
	if !totp.Validate(totpCode, key.Secret()) {
		return "", errorCode.ErrUnauthorized
	}
	return TicketString(ticketRepository.CreateTicket("TOTP", account.ID())), nil
}

func createEmailTicket(account core.Account, verificationCode string) (string, error) {
	code := verificationCodeRepository.Get(account.Email, constants.TICKET)
	if !verify(code, verificationCode) {
		return "", errorCode.ErrVerificationCode
	}
	verificationCodeRepository.Delete(account.Email, constants.TICKET)
	return TicketString(ticketRepository.CreateTicket("TOTP", account.ID())), nil
}

func createPasswordTicket(account core.Account, password string) (string, error) {
	if !utils.PasswordsMatch(account.Password, password, account.Salt) {
		return "", errorCode.ErrUnauthorized
	}
	return TicketString(ticketRepository.CreateTicket("TOTP", account.ID())), nil
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
