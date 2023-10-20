package services

import (
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model"
	"github.com/boardware-cloud/model/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

var emailSender notifications.Sender

var accountRepository core.AccountRepository
var verificationCodeRepository core.VerificationCodeRepository
var ticketRepository core.TicketRepository

func init() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
	DB, err = model.NewConnection(user, password, host, port, database)
	emailSender.SmtpHost = viper.GetString("smtp.host")
	emailSender.Port = viper.GetString("smtp.port")
	emailSender.Email = viper.GetString("smtp.email")
	emailSender.Password = viper.GetString("smtp.password")
	if err != nil {
		panic(err)
	}
	core.Init(DB)
	utils.Init()
	accountRepository = core.NewAccountRepository(DB)
	verificationCodeRepository = core.NewVerificationCodeRepository(DB)
	ticketRepository = core.NewTicketRepository(DB)
}
