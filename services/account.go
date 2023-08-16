package services

import (
	"time"

	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/errors"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model/core"
)

const EXPIRED_TIME = 60 * 5
const MAX_TRIES = 10

type Account struct {
	ID    uint           `json:"id"`
	Email string         `json:"email"`
	Role  constants.Role `json:"role"`
}

type Session struct {
	Account     Account               `json:"account"`
	Token       string                `json:"token"`
	TokeType    constants.TokenType   `json:"tokenType"`
	TokenFormat constants.TokenFormat `json:"tokenFormat"`
	ExpiredAt   int64                 `json:"expiredAt"`
	CreatedAt   int64                 `json:"createdAt"`
}

func (a Account) Forward() core.Account {
	return core.Account{
		ID:    a.ID,
		Email: a.Email,
		Role:  a.Role,
	}
}

func (a *Account) Backward(account core.Account) *Account {
	a.Email = account.Email
	a.ID = account.ID
	a.Role = account.Role
	return a
}

func CreateSession(email, password string) (*Session, *errors.Error) {
	var account *core.Account
	DB.First(&account, "email = ?", email)
	if account == nil || !utils.PasswordsMatch(account.Password, password, account.Salt) {
		return nil, errors.AuthenticationError()
	}
	expiredAt := time.Now().AddDate(0, 0, 7).Unix()
	token, err := utils.SignJwt(
		utils.UintToString(account.ID),
		account.Email,
		string(account.Role),
		expiredAt,
	)
	if err != nil {
		return nil, errors.UndefineError()
	}
	var a Account
	a.Backward(*account)
	return &Session{
		Account:     *a.Backward(*account),
		Token:       token,
		TokenFormat: constants.JWT,
		TokeType:    constants.BEARER,
		ExpiredAt:   expiredAt,
	}, nil
}

func CreateAccount(email, password string, role constants.Role) (*Account, *errors.Error) {
	var accounts []core.Account
	DB.Find(&accounts, "email = ?", email)
	if len(accounts) > 0 {
		return nil, errors.EmailExists()
	}
	hashed, salt := utils.HashWithSalt(password)
	account := Account{
		Email: email,
		Role:  role,
	}.Forward()
	account.Password = hashed
	account.Salt = salt
	if role != "" {
		account.Role = role
	} else {
		account.Role = constants.USER
	}
	DB.Create(&account)
	var back Account
	return back.Backward(account), nil
}

func GetAccountById(id uint) *Account {
	var coreAccount core.Account
	if ctx := DB.Find(&coreAccount, id); ctx.RowsAffected == 0 {
		return nil
	}
	var account Account
	return account.Backward(coreAccount)
}

func GetAccountByEmail(email string) *Account {
	var coreAccount core.Account
	if ctx := DB.Where("email = ?", email).Find(&coreAccount); ctx.RowsAffected == 0 {
		return nil
	}
	var account Account
	return account.Backward(coreAccount)
}

func CreateAccountWithVerificationCode(email, code, password string) (*Account, *errors.Error) {
	verificationCode := GetVerification(email, constants.CREATE_ACCOUNT)
	if !verify(verificationCode, code) {
		return nil, errors.VerificationCodeError()
	}
	DB.Delete(&verificationCode)
	return CreateAccount(email, password, constants.USER)
}

func SetPassword(account core.Account, newPassword string) {
	hashed, salt := utils.HashWithSalt(newPassword)
	account.Password = hashed
	account.Salt = salt
	DB.Save(&account)
}

func UpdatePassword(email string, password *string, code *string, newPassword string) *errors.Error {
	var account core.Account
	ctx := DB.Where("email = ?", email).Find(&account)
	if ctx.RowsAffected == 0 {
		return errors.NotFoundError()
	}
	if password != nil {
		if !utils.PasswordsMatch(account.Password, *password, account.Salt) {
			return errors.AuthenticationError()
		}
		SetPassword(account, newPassword)
		return nil
	}
	if code != nil {
		verificationCode := GetVerification(email, constants.SET_PASSWORD)
		if !verify(verificationCode, *code) {
			return errors.VerificationCodeError()
		}
		DB.Delete(verificationCode)
		SetPassword(account, newPassword)
		return nil
	}
	return errors.VerificationCodeError()
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
