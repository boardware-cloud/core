package main

import (
	"github.com/Dparty/common/config"
	"github.com/boardware-cloud/core/controllers"
	_ "github.com/boardware-cloud/core/services"
)

func main() {
	port := ":" + config.GetString("server.port")
	controllers.Run(port)
}
