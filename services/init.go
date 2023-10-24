package services

import (
	"github.com/boardware-cloud/common/notifications"
	coreModel "github.com/boardware-cloud/model/core"
	"gorm.io/gorm"
)

var DB *gorm.DB

var emailSender notifications.Sender

var accountRepository coreModel.AccountRepository
var verificationCodeRepository coreModel.VerificationCodeRepository
var ticketRepository coreModel.TicketRepository
var webauthRepository coreModel.WebauthRepository

func Init(db *gorm.DB) {
	DB = db
	coreModel.Init(DB)
	accountRepository = coreModel.NewAccountRepository(DB)
	verificationCodeRepository = coreModel.NewVerificationCodeRepository(DB)
	ticketRepository = coreModel.NewTicketRepository(DB)
	webauthRepository = coreModel.NewWebauthRepository(DB)
}
