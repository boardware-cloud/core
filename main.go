package main

import (
	"github.com/Dparty/common/config"
	"github.com/boardware-cloud/core/controllers"
	_ "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/model"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	port := ":" + config.GetString("server.port")
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	dbport := config.GetString("database.port")
	database := config.GetString("database.database")
	db, _ = model.NewConnection(user, password, host, dbport, database)
	controllers.Init(db)
	controllers.Run(port)
}
