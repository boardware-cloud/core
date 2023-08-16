package controllers

import (
	api "github.com/boardware-cloud/core-api"
	core "github.com/boardware-cloud/core/services"

	"github.com/boardware-cloud/common/errors"

	"github.com/boardware-cloud/common/constants"

	"net/http"

	"github.com/boardware-cloud/common/utils"

	"github.com/gin-gonic/gin"
)

type AccountApi struct{}

var accountApi AccountApi

func (AccountApi) CreateSession(c *gin.Context, createSessionRequest api.CreateSessionRequest) {
	session, sessionError := core.CreateSession(
		createSessionRequest.Email,
		createSessionRequest.Password,
	)
	if sessionError != nil {
		sessionError.GinHandler(c)
		return
	}
	c.JSON(http.StatusCreated, SessionBackward(*session))
}

func (AccountApi) CreateAccount(c *gin.Context, createAccountRequest api.CreateAccountRequest) {
	if createAccountRequest.VerificationCode != nil {
		a, createError := core.CreateAccountWithVerificationCode(
			createAccountRequest.Email,
			*createAccountRequest.VerificationCode,
			createAccountRequest.Password)
		if createError != nil {
			createError.GinHandler(c)
			return
		}
		c.JSON(http.StatusCreated, AccountBackward(*a))
		return
	}
	middleware.IsRoot(c, func(_ *gin.Context, _ api.Account) {
		var createAccountRequest api.CreateAccountRequest
		err := c.ShouldBindJSON(&createAccountRequest)
		if err != nil {
			// TODO: Error message
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		role := constants.USER
		if createAccountRequest.Role != nil {
			role = constants.Role(*createAccountRequest.Role)
		}
		a, createError := core.CreateAccount(
			createAccountRequest.Email,
			createAccountRequest.Password,
			role,
		)
		if createError != nil {
			createError.GinHandler(c)
			return
		}
		c.JSON(http.StatusCreated, AccountBackward(*a))
	})
}

func (AccountApi) ListAccount(gin_context *gin.Context, ordering api.Ordering, index int64, limit int64) {
	// TODO: List account api
}

func (AccountApi) GetAccount(c *gin.Context) {
	auth := Authorize(c)
	if auth.Status != Authorized {
		errors.UnauthorizedError().GinHandler(c)
		return
	}
	account := core.GetAccountById(auth.AccountId)
	if account == nil {
		errors.NotFoundError().GinHandler(c)
		return
	}
	c.JSON(http.StatusOK, AccountBackward(*account))
}

func (AccountApi) GetAccountById(id string) *api.Account {
	account := core.GetAccountById(utils.StringToUint(id))
	if account == nil {
		return nil
	}
	a := AccountBackward(*account)
	return &a
}

func (a AccountApi) VerifySession(c *gin.Context, sessionVerificationRequest api.SessionVerificationRequest) {
	auth := Authorize(c)
	if auth.Status != Authorized {
		errors.UnauthorizedError().GinHandler(c)
		return
	}
	account := a.GetAccountById(utils.UintToString(auth.AccountId))
	if account == nil {
		c.JSON(401, gin.H{})
		return
	}
	c.JSON(http.StatusOK, api.Session{
		Account: *account,
	})
}

func (AccountApi) UpdatePassword(c *gin.Context, request api.UpdatePasswordRequest) {
	err := core.UpdatePassword(request.Email, request.Password, request.VerificationCode, request.NewPassword)
	if err != nil {
		err.GinHandler(c)
		return
	}
	c.String(http.StatusNoContent, "")
}
