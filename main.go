package main

import (
	"github.com/boardware-cloud/core/controllers"
	_ "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	config := viper.New()
	config.SetConfigName("env")
	config.SetConfigType("yaml")
	config.AddConfigPath("./config")
	config.ReadInConfig()
	port := ":" + config.GetString("server.port")
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	dbport := config.GetString("database.port")
	database := config.GetString("database.database")
	db, _ = model.NewConnection(user, password, host, dbport, database)
	controllers.Init(config, db)
	controllers.Run(port)
}
