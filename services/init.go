package services

import (
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/model"
	"github.com/boardware-cloud/model/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

var emailSender notifications.Sender

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
