package main

import (
	"fmt"

	"github.com/boardware-cloud/core/controllers"
	_ "github.com/boardware-cloud/core/services"

	"github.com/spf13/viper"
)

func main() {
	var err error
	viper.SetConfigName(".env") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")   // optionally look for config in the working directory
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	port := ":" + viper.GetString("server.port")
	controllers.Init()
	controllers.Run(port)
}
