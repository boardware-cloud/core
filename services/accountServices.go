package services

import (
	"github.com/boardware-cloud/core/services/model"
	"github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAccountService(inject *gorm.DB) AccountService {
	return AccountService{accountRepository: core.NewAccountRepository(inject)}
}

type AccountService struct {
	accountRepository core.AccountRepository
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

func (a AccountService) GetAccount(id uint) *model.Account {
	account := a.accountRepository.GetById(id)
	if account == nil {
		return nil
	}
	return &model.Account{Entity: *account}
}
