package services

import (
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/model"
	"github.com/boardware-cloud/model/core"
	"gorm.io/gorm"
)

var DB *gorm.DB

var emailSender notifications.Sender

func init() {
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetString("database.port")
	database := config.GetString("database.database")
	var err error
	DB, err = model.NewConnection(user, password, host, port, database)
	emailSender.SmtpHost = config.GetString("smtp.host")
	emailSender.Port = config.GetString("smtp.port")
	emailSender.Email = config.GetString("smtp.email")
	emailSender.Password = config.GetString("smtp.password")
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&core.Account{})
	DB.AutoMigrate(&core.Service{})
	DB.AutoMigrate(&core.VerificationCode{})
}
