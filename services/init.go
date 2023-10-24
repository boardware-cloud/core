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

func Init(db *gorm.DB) {
	DB = db
	coreModel.Init(DB)
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	emailSender.SmtpHost = viper.GetString("smtp.host")
	emailSender.Port = viper.GetString("smtp.port")
	emailSender.Email = viper.GetString("smtp.email")
	emailSender.Password = viper.GetString("smtp.password")
	accountRepository = coreModel.NewAccountRepository(DB)
	verificationCodeRepository = coreModel.NewVerificationCodeRepository(DB)
	ticketRepository = coreModel.NewTicketRepository(DB)
	webauthRepository = coreModel.NewWebauthRepository(DB)
	ticketService = NewTicketService(db)
}
