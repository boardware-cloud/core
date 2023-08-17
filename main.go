package main

import (
	"errors"
	"log"
	"net/smtp"

	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/core/controllers"
	_ "github.com/boardware-cloud/core/services"
)

type loginAuth struct {
	username, password string
}

// LoginAuth is used for smtp login auth
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func main() {
	user := config.GetString("smtp.email")
	from := config.GetString("smtp.email")
	password := config.GetString("smtp.password")
	to := []string{"chenyunda218@gmail.com"}
	smtpHost := config.GetString("smtp.host")
	smtpPort := config.GetString("smtp.port")

	message := []byte("Hello! I'm trying out smtp to send emails to recipients.")

	// Create authentication
	auth := LoginAuth(user, password)

	// Send actual message
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	if err != nil {
		log.Fatal(err)
	}
	port := ":" + config.GetString("server.port")
	controllers.Init()
	controllers.Run(port)
}
