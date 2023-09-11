package controllers

import (
	"net/http"

	errorCode "github.com/boardware-cloud/common/code"
	"github.com/boardware-cloud/common/constants"
	api "github.com/boardware-cloud/core-api"
	"github.com/boardware-cloud/core/services"
	"github.com/gin-gonic/gin"
)

type VerificationApi struct{}

var verificationApi VerificationApi

const CREATE_INTERVAL = 60

func (VerificationApi) CreateVerificationCode(c *gin.Context, request api.CreateVerificationCodeRequest) {
	purpose := constants.VerificationCodePurpose(request.Purpose)
	if request.Email == nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	err := services.CreateVerificationCode(*request.Email, purpose)
	if err != nil {
		errorCode.GinHandler(c, err)
		return
	}
	c.JSON(http.StatusCreated, api.CreateVerificationCodeRespones{
		Email:   request.Email,
		Purpose: request.Purpose,
		Result:  api.SUCCESS_CREATED,
	})
}
