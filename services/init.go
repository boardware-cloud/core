package services

import (
	"fmt"

	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/model"
	"github.com/boardware-cloud/model/core"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

var emailSender notifications.Sender

func init() {
	var err error
	viper.SetConfigName(".env") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")   // optionally look for config in the working directory
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("databases fatal error config file: %w", err))
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
	DB.AutoMigrate(&core.Account{})
	DB.AutoMigrate(&core.Service{})
	DB.AutoMigrate(&core.VerificationCode{})
}
