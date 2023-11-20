package services

import (
	"github.com/boardware-cloud/common/notifications"
	coreModel "github.com/boardware-cloud/model/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

var emailSender notifications.Sender

var accountRepository coreModel.AccountRepository
var verificationCodeRepository coreModel.VerificationCodeRepository
var ticketRepository coreModel.TicketRepository
var webauthRepository coreModel.WebauthRepository
var ticketService TicketService

func Init(config *viper.Viper, db *gorm.DB) {
	DB = db
	coreModel.Init(DB)
	emailSender.SmtpHost = config.GetString("smtp.host")
	emailSender.Port = config.GetString("smtp.port")
	emailSender.Email = config.GetString("smtp.email")
	emailSender.Password = config.GetString("smtp.password")
	accountRepository = coreModel.NewAccountRepository(DB)
	verificationCodeRepository = coreModel.NewVerificationCodeRepository(DB)
	ticketRepository = coreModel.NewTicketRepository(DB)
	webauthRepository = coreModel.NewWebauthRepository(DB)
	ticketService = NewTicketService(db)
}
