package services

import (
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/model"
	coreModel "github.com/boardware-cloud/model/core"
	"gorm.io/gorm"
)

var DB *gorm.DB

var emailSender notifications.Sender

var accountRepository = coreModel.GetAccountRepository()
var verificationCodeRepository = coreModel.GetVerificationCodeRepository()
var ticketRepository = coreModel.GetTicketRepository()
var webauthRepository = coreModel.GetWebauthRepository()

func init() {
	DB = model.GetDB()
	emailSender.SmtpHost = config.GetString("smtp.host")
	emailSender.Port = config.GetString("smtp.port")
	emailSender.Email = config.GetString("smtp.email")
	emailSender.Password = config.GetString("smtp.password")
}
