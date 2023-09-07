package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/boardware-cloud/common/code"
	"github.com/boardware-cloud/common/constants/authenication"
	"github.com/boardware-cloud/common/errors"
	api "github.com/boardware-cloud/core-api"
	core "github.com/boardware-cloud/core/services"
	"github.com/chenyunda218/golambda"
	"github.com/go-webauthn/webauthn/protocol"

	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/middleware"
	model "github.com/boardware-cloud/model/core"

	"github.com/gin-gonic/gin"
)

type AccountApi struct{}

// DeleteTotp implements coreapi.AccountApiInterface.
func (AccountApi) DeleteTotp(ctx *gin.Context) {
	middleware.GetAccount(ctx,
		func(ctx *gin.Context, account model.Account) {
			core.DeleteTotp(account)
			ctx.JSON(http.StatusNoContent, "")
		})
}

// GetAuthentication implements coreapi.AccountApiInterface.
func (AccountApi) GetAuthentication(ctx *gin.Context, email string) {
	if email == "" {
		errors.NotFoundError().GinHandler(ctx)
		return
	}
	factors := core.GetAuthenticationFactors(email)
	if len(factors) == 0 {
		errors.NotFoundError().GinHandler(ctx)
		return
	}
	ctx.JSON(http.StatusOK, api.Authentication{Factors: golambda.Map(factors,
		func(_ int, factor authenication.AuthenticationFactor) string {
			return string(factor)
		})})
}

// DeleteWebAuthn implements coreapi.AccountApiInterface.
func (AccountApi) DeleteWebAuthn(ctx *gin.Context, id string) {
	middleware.GetAccount(ctx, func(c *gin.Context, account model.Account) {
		err := core.DeleteWebAuthn(account, utils.StringToUint(id))
		if err != nil {
			err.GinHandler(ctx)
			return
		}
		ctx.JSON(http.StatusNoContent, "")
	})
}

// CreateWebauthnTickets implements coreapi.AccountApiInterface.
func (AccountApi) CreateWebauthnTickets(ctx *gin.Context, id string) {
	response, err := protocol.ParseCredentialRequestResponseBody(ctx.Request.Body)
	if err != nil {
		return
	}
	ticket, err := core.FinishLogin(utils.StringToUint(id), response)
	if err != nil {
		code.GinHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, api.Ticket{
		Token: ticket,
		Type:  api.WEBAUTHN,
	})
}

// ListWebAuthn implements coreapi.AccountApiInterface.
func (AccountApi) ListWebAuthn(ctx *gin.Context) {
	middleware.GetAccount(ctx, func(ctx *gin.Context, account model.Account) {
		webauthns := core.ListWebAuthn(account)
		ctx.JSON(http.StatusOK, golambda.Map(webauthns,
			func(_ int, cred model.Credential) api.WebAuthn {
				return api.WebAuthn{
					Id:        utils.UintToString(cred.ID),
					Name:      cred.Name,
					Os:        cred.Os,
					CreatedAt: cred.CreatedAt.Unix(),
				}
			}))
	})
}

// CreateWebauthnTicketChallenge implements coreapi.AccountApiInterface.
func (AccountApi) CreateWebauthnTicketChallenge(ctx *gin.Context, request api.CreateTicketChallenge) {
	account, err := model.GetAccountByEmail(request.Email)
	if err != nil {
		code.GinHandler(ctx, err)
		return
	}
	option, session, err := core.BeginLogin(account)
	if err != nil {
		code.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":        utils.UintToString(session.ID),
		"publicKey": option.Response,
	})
}

// CreateWebAuthnChallenge implements coreapi.AccountApiInterface.
func (AccountApi) CreateWebAuthnChallenge(ctx *gin.Context) {
	middleware.GetAccount(ctx,
		func(ctx *gin.Context, account model.Account) {
			options, session := core.BeginRegistration(account)
			ctx.JSON(http.StatusOK, gin.H{
				"id":        utils.UintToString(session.ID),
				"publicKey": options.Response,
			})
		})
}

type Credential struct {
	protocol.CredentialCreationResponse
	Name string `json:"name"`
	Os   string `json:"os"`
}

// CreateWebauthn implements coreapi.AccountApiInterface.
func (AccountApi) CreateWebauthn(ctx *gin.Context, id string) {
	middleware.GetAccount(ctx,
		func(ctx *gin.Context, account model.Account) {
			var ccr Credential
			if err := json.NewDecoder(ctx.Copy().Request.Body).Decode(&ccr); err != nil {
				ctx.JSON(http.StatusBadRequest, "")
				return
			}
			if err := core.FinishRegistration(
				account,
				utils.StringToUint(id),
				ccr.Name,
				ccr.Os,
				ccr.CredentialCreationResponse); err != nil {
				err.GinHandler(ctx)
				return
			}
			ctx.JSON(http.StatusCreated, "")
		})
}

// GetTotp implements coreapi.AccountApiInterface.
func (AccountApi) GetTotp(c *gin.Context) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			c.JSON(http.StatusOK, api.Totp{Url: core.CreateTotp(account)})
		})
}

// CreateTotp2FA implements coreapi.AccountApiInterface.
func (AccountApi) CreateTotp2FA(c *gin.Context, request api.PutTotpRequest) {
	middleware.GetAccount(c, func(c *gin.Context, account model.Account) {
		err := core.NFactor(account, request.Tickets, 1)
		if err != nil {
			err.GinHandler(c)
			return
		}
		url, err := core.UpdateTotp2FA(account, request.Url, request.TotpCode)
		if err != nil {
			err.GinHandler(c)
			return
		}
		c.JSON(http.StatusOK, api.Totp{
			Url: url,
		})
	})
}

func (AccountApi) CreateSession(c *gin.Context, createSessionRequest api.CreateSessionRequest) {
	if createSessionRequest.Tickets != nil {
		session, sessionError := core.CreateSessionWithTickets(*createSessionRequest.Email, *createSessionRequest.Tickets)
		if sessionError != nil {
			sessionError.GinHandler(c)
			return
		}
		c.JSON(http.StatusCreated, SessionBackward(*session))
	} else {
		session, sessionError := core.CreateSession(
			*createSessionRequest.Email,
			createSessionRequest.Password,
			createSessionRequest.VerificationCode,
			createSessionRequest.TotpCode,
		)
		if sessionError != nil {
			sessionError.GinHandler(c)
			return
		}
		c.JSON(http.StatusCreated, SessionBackward(*session))
	}
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
	middleware.IsRoot(c, func(_ *gin.Context, _ model.Account) {
		var createAccountRequest api.CreateAccountRequest
		err := c.ShouldBindJSON(&createAccountRequest)
		if err != nil {
			// TODO: Error message
			c.JSON(http.StatusBadRequest, "")
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

func (AccountApi) ListAccount(ctx *gin.Context, ordering api.Ordering, index int64, limit int64) {
	middleware.IsRoot(ctx,
		func(ctx *gin.Context, account model.Account) {
			list := core.ListAccount(index, limit)
			ctx.JSON(http.StatusOK, api.AccountList{
				Data: golambda.Map(list.Data, func(_ int, account core.Account) api.Account {
					return api.Account{
						Id:      utils.UintToString(account.ID),
						Email:   account.Email,
						HasTotp: account.HasTotp,
						Role:    api.Role(account.Role),
					}
				}),
				Pagination: api.Pagination{
					Limit: list.Pagination.Limit,
					Index: list.Pagination.Index,
					Total: list.Pagination.Total,
				},
			})
		})
}

func (AccountApi) GetAccount(ctx *gin.Context) {
	middleware.GetAccount(ctx,
		func(c *gin.Context, a model.Account) {
			account := core.GetAccountById(a.ID)
			c.JSON(http.StatusOK, AccountBackward(*account))
		})
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
	auth := middleware.Authorize(c)
	if auth.Status != middleware.Authorized {
		errors.UnauthorizedError().GinHandler(c)
		return
	}
	account := a.GetAccountById(utils.UintToString(auth.AccountId))
	if account == nil {
		c.JSON(http.StatusUnauthorized, "")
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

var accountApi AccountApi
