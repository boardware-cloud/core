package services

import (
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model/core"
	"gorm.io/gorm"
)

var DB *gorm.DB

var emailSender notifications.Sender

var accountRepository core.AccountRepository
var verificationCodeRepository core.VerificationCodeRepository
var ticketRepository core.TicketRepository
var webauthRepository core.WebauthRepository

func Init(db *gorm.DB) {
	core.Init(DB)
	utils.Init()
	accountRepository = core.NewAccountRepository(DB)
	verificationCodeRepository = core.NewVerificationCodeRepository(DB)
	ticketRepository = core.NewTicketRepository(DB)
	webauthRepository = core.NewWebauthRepository(DB)
}
