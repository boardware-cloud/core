package services

import (
	"time"

	"github.com/Dparty/common/fault"
	errorCode "github.com/boardware-cloud/common/code"
	constants "github.com/boardware-cloud/common/constants/account"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var accountService *AccountService

func GetAccountService() *AccountService {
	if accountService == nil {
		accountService = NewAccountService()
	}
	return accountService
}

func NewAccountService() *AccountService {
	return &AccountService{accountRepository: core.GetAccountRepository()}
}

type AccountService struct {
	accountRepository *core.AccountRepository
}

func (a AccountService) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := Authorize(c)
		if auth.Status == Authorized {
			account := a.GetAccount(auth.AccountId)
			c.Set("account", *account)
		}
		c.Next()
	}
}

func (a AccountService) GetAccount(id uint) *Account {
	account := a.accountRepository.GetById(id)
	if account == nil {
		return nil
	}
	return &Account{Entity: *account}
}
func (a AccountService) GetAccountByEmail(email string) *Account {
	account := a.accountRepository.GetByEmail(email)
	if account == nil {
		return nil
	}
	return &Account{Entity: *account}
}

func (a AccountService) GetByEmail(email string) *Account {
	var account *core.Account = a.accountRepository.GetByEmail(email)
	if account == nil {
		return nil
	}
	return &Account{Entity: *account}
}

func (a AccountService) CreateAccount(email, password string, role constants.Role) (*Account, error) {
	account, err := a.accountRepository.Create(email, password, role)
	if err != nil {
		return nil, err
	}
	return &Account{Entity: *account}, nil
}

func (a AccountService) CreateAccountWithVerificationCode(email, code, password string) (*Account, error) {
	if email == "" {
		return nil, errorCode.ErrBadRequest
	}
	verificationCode := verificationCodeRepository.Get(email, constants.CREATE_ACCOUNT)
	if !verify(verificationCode, code) {
		return nil, errorCode.ErrVerificationCode
	}
	DB.Delete(&verificationCode)
	return a.CreateAccount(email, password, constants.USER)
}

func (a AccountService) UpdatePassword(email string, password *string, code *string, newPassword string) error {
	var ac = a.GetAccountByEmail(email)
	if ac == nil {
		return errorCode.ErrNotFound
	}
	var account = ac.Entity
	if password != nil {
		if !utils.PasswordsMatch(account.Password, *password, account.Salt) {
			return errorCode.ErrUnauthorized
		}
		ac.UpdatePassword(newPassword)
		return nil
	}
	if code != nil {
		verificationCode := verificationCodeRepository.Get(email, constants.SET_PASSWORD)
		if !verify(verificationCode, *code) {
			return errorCode.ErrVerificationCode
		}
		verificationCodeRepository.Delete(email, constants.SET_PASSWORD)
		ac.UpdatePassword(newPassword)
		return nil
	}
	return errorCode.ErrVerificationCode
}

func (a AccountService) UpdateTotp2FA(account core.Account, url, totpCode string) (string, error) {
	key, err := otp.NewKeyFromURL(url)
	if err != nil {
		return "", errorCode.ErrUnauthorized
	}
	if !totp.Validate(totpCode, key.Secret()) {
		return "", errorCode.ErrUnauthorized
	}
	account.Totp = &url
	a.accountRepository.Save(&account)
	return *account.Totp, nil
}

const CREATE_SESSION_COLDDOWN = 500

func (a AccountService) CreateSessionWithTickets(email string, tokens []string) (*Session, error) {
	account := a.accountRepository.GetByEmail(email)
	if account == nil {
		return nil, fault.ErrUnauthorized
	}
	if !account.ColdDown(CREATE_SESSION_COLDDOWN) {
		return nil, errorCode.ErrTooManyRequests
	}
	account.CreateColdDown()
	err := NFactor(*account, tokens, 2)
	if err != nil {
		return nil, errorCode.ErrUnauthorized
	}
	expiredAt := time.Now().AddDate(0, 0, 7).Unix()
	token, _ := utils.SignJwt(
		utils.UintToString(account.ID()),
		account.Email,
		string(account.Role),
		expiredAt,
		"ACTIVED",
	)
	return &Session{
		Token:       token,
		TokenFormat: constants.JWT,
		TokeType:    constants.BEARER,
		ExpiredAt:   expiredAt,
	}, nil
}

func CheckIfRoot(ctx *gin.Context) *Account {
	accountInterface, ok := ctx.Get("account")
	if !ok {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return nil
	}
	account, ok := accountInterface.(Account)
	if !ok {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return nil
	}
	if account.Role() != "ROOT" {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return nil
	}
	return &account
}

func GetAccount(ctx *gin.Context) *Account {
	accountInterface, ok := ctx.Get("account")
	if !ok {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return nil
	}
	account, ok := accountInterface.(Account)
	if !ok {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return nil
	}
	return &account
}
