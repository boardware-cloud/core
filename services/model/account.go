package model

import (
	"time"

	constants "github.com/boardware-cloud/common/constants/account"
	"github.com/boardware-cloud/model/core"
)

type Account struct {
	Entity core.Account `json:"entity"`
}

func (a Account) ID() uint {
	return a.Entity.ID
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
