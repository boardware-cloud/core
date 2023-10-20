package services

import (
	"strings"

	"github.com/boardware-cloud/common/utils"
	"github.com/gin-gonic/gin"
)

type Headers struct {
	Authorization string
}

type Authentication struct {
	Status    AuthenticationStatus
	Role      string
	AccountId uint
	Email     string
}

type AuthenticationStatus string

const (
	Authorized   AuthenticationStatus = "Authorized"
	Unauthorized AuthenticationStatus = "Unauthorized"
)

func Authorize(c *gin.Context) Authentication {
	var headers Headers
	c.ShouldBindHeader(&headers)
	authorization := headers.Authorization
	splited := strings.Split(authorization, " ")
	if authorization == "" || len(splited) != 2 {
		return Authentication{
			Status: Unauthorized,
		}
	}
	return AuthorizeByJWT(splited[1])
}

func AuthorizeByJWT(token string) Authentication {
	claims, err := utils.VerifyJwt(token)
	if err != nil {
		return Authentication{
			Status: Unauthorized,
		}
	}
	return Authentication{
		Status:    Authorized,
		Email:     claims["email"].(string),
		AccountId: utils.StringToUint(claims["id"].(string)),
		Role:      claims["role"].(string),
	}
}
