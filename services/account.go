package services

import (
	"time"

	errorCode "github.com/boardware-cloud/common/code"
	constants "github.com/boardware-cloud/common/constants/account"
	"github.com/boardware-cloud/common/constants/authenication"
	"github.com/boardware-cloud/model/core"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const EXPIRED_TIME = 60 * 5
const MAX_TRIES = 10

type Account struct {
	Entity core.Account `json:"entity"`
}

func (a Account) ID() uint {
	return a.Entity.ID()
}

func (a Account) Email() string {
	return a.Entity.Email
}

func (a Account) Role() constants.Role {
	return a.Entity.Role
}

func (a Account) HasTotp() bool {
	return a.Entity.Totp != nil
}

func (a Account) RegisteredOn() time.Time {
	return a.Entity.CreatedAt
}

func (a Account) ListWebAuthn() []core.Credential {
	return webauthRepository.List("account_id = ?", a.ID())
}

func (a *Account) DeleteTotp() *Account {
	a.Entity.Totp = nil
	a.Submit()
	return a
}

func (a *Account) DeleteWebAuthn(id uint) *Account {
	DB.Where("account_id = ? AND id = ?", a.ID(), id).Delete(&core.Credential{})
	return a
}

func (a *Account) UpdatePassword(newPassword string) *Account {
	a.Entity.SetPassword(newPassword)
	a.Submit()
	return a
}

func (a *Account) CreateTotp() string {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "cloud.boardware.com",
		AccountName: a.Email(),
	})
	return key.String()
}

func (a *Account) UpdateTotp2FA(url, totpCode string) (string, error) {
	key, err := otp.NewKeyFromURL(url)
	if err != nil {
		return "", errorCode.ErrUnauthorized
	}
	if !totp.Validate(totpCode, key.Secret()) {
		return "", errorCode.ErrUnauthorized
	}
	a.Entity.Totp = &url
	a.Submit()
	return *a.Entity.Totp, nil
}

func (a *Account) Submit() *Account {
	accountRepository.Save(&a.Entity)
	return a
}

type Session struct {
	Token       string                `json:"token"`
	TokeType    constants.TokenType   `json:"tokenType"`
	TokenFormat constants.TokenFormat `json:"tokenFormat"`
	ExpiredAt   int64                 `json:"expiredAt"`
	CreatedAt   int64                 `json:"createdAt"`
	Status      string                `json:"status"`
}

func GetAuthenticationFactors(email string) []authenication.AuthenticationFactor {
	account := accountRepository.GetByEmail(email)
	var factors []authenication.AuthenticationFactor
	if account == nil {
		return factors
	}
	factors = append(factors, authenication.PASSWORD)
	factors = append(factors, authenication.EMAIL)
	if account.Totp != nil {
		factors = append(factors, authenication.TOTP)
	}
	if len(account.WebAuthnCredentials()) != 0 {
		factors = append(factors, authenication.WEBAUTHN)
	}
	return factors
}

func UpdateUserRole(accountId, role constants.Role) {
	// TODO:
}

func NFactor(account core.Account, tokens []string, factor int) error {
	var fa map[string]bool = make(map[string]bool)
	for _, token := range tokens {
		ticket, err := ticketService.Get().UseTicket(token)
		if err != nil {
			return errorCode.ErrUnauthorized
		}
		if ticket.AccountId == account.ID() {
			fa[ticket.Type] = true
		}
	}
	if len(fa) < factor {
		return errorCode.ErrUnauthorized
	}
	return nil
}

func verify(v *core.VerificationCode, code string) bool {
	if v == nil {
		return false
	}
	v.Tries++
	verificationCodeRepository.Save(v)
	if v.Code != code || time.Now().Unix()-v.CreatedAt.Unix() > EXPIRED_TIME || v.Tries > MAX_TRIES {
		return false
	}
	return true
}
