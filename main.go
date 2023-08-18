package main

import (
	"github.com/boardware-cloud/core/controllers"
	_ "github.com/boardware-cloud/core/services"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.ReadInConfig()
	port := ":" + viper.GetString("server.port")
	controllers.Init()
	controllers.Run(port)
}
