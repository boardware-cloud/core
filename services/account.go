package services

import (
	"time"

	errorCode "github.com/boardware-cloud/common/code"
	constants "github.com/boardware-cloud/common/constants/account"
	"github.com/boardware-cloud/common/constants/authenication"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/core"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const EXPIRED_TIME = 60 * 5
const MAX_TRIES = 10

type Account struct {
	ID           uint           `json:"id"`
	Email        string         `json:"email"`
	Role         constants.Role `json:"role"`
	HasTotp      bool           `json:"hasTotp"`
	RegisteredOn time.Time      `json:"registeredOn"`
}

type Session struct {
	Token       string                `json:"token"`
	TokeType    constants.TokenType   `json:"tokenType"`
	TokenFormat constants.TokenFormat `json:"tokenFormat"`
	ExpiredAt   int64                 `json:"expiredAt"`
	CreatedAt   int64                 `json:"createdAt"`
	Status      string                `json:"status"`
}

func UpdateTotp2FA(account core.Account, url, totpCode string) (string, error) {
	key, err := otp.NewKeyFromURL(url)
	if err != nil {
		return "", errorCode.ErrUnauthorized
	}
	if !totp.Validate(totpCode, key.Secret()) {
		return "", errorCode.ErrUnauthorized
	}
	account.Totp = &url
	DB.Save(&account)
	return *account.Totp, nil
}

func CreateTotp(account core.Account) string {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "cloud.boardware.com",
		AccountName: account.Email,
	})
	return key.String()
}

func DeleteTotp(account core.Account) {
	account.Totp = nil
	DB.Save(&account)
}

func CreateSessionWithTickets(email string, tokens []string) (*Session, error) {
	account, err := core.FindAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	if !account.ColdDown(500) {
		return nil, errorCode.ErrTooManyRequests
	}
	account.CreateColdDown()
	err = NFactor(account, tokens, 2)
	if err != nil {
		return nil, errorCode.ErrUnauthorized
	}
	expiredAt := time.Now().AddDate(0, 0, 7).Unix()
	token, _ := utils.SignJwt(
		utils.UintToString(account.ID),
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

func GetAuthenticationFactors(email string) []authenication.AuthenticationFactor {
	account, err := core.FindAccountByEmail(email)
	var factors []authenication.AuthenticationFactor
	if err != nil {
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

func CreateAccount(email, password string, role constants.Role) (*Account, error) {
	if ctx := DB.Find(&core.Account{}, "email = ?", email); ctx.RowsAffected != 0 {
		return nil, errorCode.ErrEmailExists
	}
	hashed, salt := utils.HashWithSalt(password)
	account := AccountForward(Account{Email: email, Role: role})
	account.Password = hashed
	account.Salt = salt
	if role != "" {
		account.Role = role
	} else {
		account.Role = constants.USER
	}
	DB.Create(&account)
	DB.Delete(&core.VerificationCode{Identity: email})
	back := AccountBackward(account)
	return &back, nil
}

func GetAccountById(id uint) *Account {
	var coreAccount core.Account
	if ctx := DB.Find(&coreAccount, id); ctx.RowsAffected == 0 {
		return nil
	}
	account := AccountBackward(coreAccount)
	return &account
}

func GetAccountByEmail(email string) *Account {
	var coreAccount core.Account
	if ctx := DB.Where("email = ?", email).Find(&coreAccount); ctx.RowsAffected == 0 {
		return nil
	}
	account := AccountBackward(coreAccount)
	return &account
}

func CreateAccountWithVerificationCode(email, code, password string) (*Account, error) {
	if email == "" {
		return nil, errorCode.ErrBadRequest
	}
	verificationCode := GetVerification(email, constants.CREATE_ACCOUNT)
	if !verify(verificationCode, code) {
		return nil, errorCode.ErrVerificationCode
	}
	DB.Delete(&verificationCode)
	return CreateAccount(email, password, constants.USER)
}

func setPassword(account core.Account, newPassword string) {
	hashed, salt := utils.HashWithSalt(newPassword)
	account.Password = hashed
	account.Salt = salt
	DB.Save(&account)
}

func UpdatePassword(email string, password *string, code *string, newPassword string) error {
	var account core.Account
	ctx := DB.Where("email = ?", email).Find(&account)
	if ctx.RowsAffected == 0 {
		return errorCode.ErrNotFound
	}
	if password != nil {
		if !utils.PasswordsMatch(account.Password, *password, account.Salt) {
			return errorCode.ErrUnauthorized
		}
		setPassword(account, newPassword)
		return nil
	}
	if code != nil {
		verificationCode := GetVerification(email, constants.SET_PASSWORD)
		if !verify(verificationCode, *code) {
			return errorCode.ErrVerificationCode
		}
		DB.Delete(verificationCode)
		setPassword(account, newPassword)
		return nil
	}
	return errorCode.ErrVerificationCode
}

func UpdateUserRole(accountId, role constants.Role) {
	// TODO:
}

func ListAccount(index, limit int64) common.List[Account] {
	return AccountListBackward(core.ListAccount(index, limit))
}

func verify(v *core.VerificationCode, code string) bool {
	if v == nil {
		return false
	}
	v.Tries++
	DB.Save(v)
	if v.Code != code || time.Now().Unix()-v.CreatedAt.Unix() > EXPIRED_TIME || v.Tries > MAX_TRIES {
		return false
	}
	return true
}

func NFactor(account core.Account, tokens []string, factor int) error {
	var fa map[string]bool = make(map[string]bool)
	for _, token := range tokens {
		ticket, err := UseTicket(token)
		if err != nil {
			return errorCode.ErrUnauthorized
		}
		if ticket.AccountId == account.ID {
			fa[ticket.Type] = true
		}
	}
	if len(fa) < factor {
		return errorCode.ErrUnauthorized
	}
	return nil
}
