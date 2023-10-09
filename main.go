package main

import (
	"context"

	"github.com/boardware-cloud/core/controllers"
	_ "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.ReadInConfig()
	port := ":" + viper.GetString("server.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	dbport := viper.GetString("database.port")
	database := viper.GetString("database.database")
	db, _ = model.NewConnection(user, password, host, dbport, database)
	ctx := context.WithValue(context.Background(), "db", db)
	controllers.Init(ctx)
	controllers.Run(port)
}
