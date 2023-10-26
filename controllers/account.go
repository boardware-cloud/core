package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Dparty/common/fault"
	errorCode "github.com/boardware-cloud/common/code"
	constants "github.com/boardware-cloud/common/constants/account"
	"github.com/boardware-cloud/common/constants/authenication"
	"github.com/boardware-cloud/common/utils"
	api "github.com/boardware-cloud/core-api"
	coreServices "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/middleware"
	model "github.com/boardware-cloud/model/core"
	"github.com/chenyunda218/golambda"
	"github.com/go-webauthn/webauthn/protocol"

	"github.com/gin-gonic/gin"
)

type AccountApi struct{}

// UpdateUserRole implements coreapi.AccountApiInterface.
func (AccountApi) UpdateUserRole(gin_context *gin.Context, id string, gin_body api.UpdateRoleRequest) {
	panic("unimplemented")
}

// ListSession implements coreapi.AccountApiInterface.
func (AccountApi) ListSession(gin_context *gin.Context) {
	panic("unimplemented")
}

func (AccountApi) GetAccountById(ctx *gin.Context, id string) {
	account := coreServices.CheckIfRoot(ctx)
	if account == nil {
		return
	}
	a := GetAccountById(id)
	if a == nil {
		errorCode.GinHandler(ctx, errorCode.ErrNotFound)
		return
	}
	ctx.JSON(http.StatusOK, a)
}

// DeleteTotp implements coreapi.AccountApiInterface.
func (AccountApi) DeleteTotp(ctx *gin.Context) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	account.DeleteTotp()
	ctx.JSON(http.StatusNoContent, "")
}

// GetAuthentication implements coreapi.AccountApiInterface.
func (AccountApi) GetAuthentication(ctx *gin.Context, email string) {
	if email == "" {
		errorCode.GinHandler(ctx, errorCode.ErrNotFound)
		return
	}
	factors := coreServices.GetAuthenticationFactors(email)
	if len(factors) == 0 {
		errorCode.GinHandler(ctx, errorCode.ErrNotFound)
		return
	}
	ctx.JSON(http.StatusOK, api.Authentication{Factors: golambda.Map(factors,
		func(_ int, factor authenication.AuthenticationFactor) string {
			return string(factor)
		})})
}

// DeleteWebAuthn implements coreapi.AccountApiInterface.
func (AccountApi) DeleteWebAuthn(ctx *gin.Context, id string) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	account.DeleteWebAuthn(utils.StringToUint(id))
	ctx.JSON(http.StatusNoContent, "")
}

// CreateWebauthnTickets implements coreapi.AccountApiInterface.
func (AccountApi) CreateWebauthnTickets(ctx *gin.Context, id string) {
	response, err := protocol.ParseCredentialRequestResponseBody(ctx.Request.Body)
	if err != nil {
		return
	}
	ticket, err := coreServices.FinishLogin(utils.StringToUint(id), response)
	if err != nil {
		errorCode.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, api.Ticket{
		Token: ticket,
		Type:  api.WEBAUTHN,
	})
}

// ListWebAuthn implements coreapi.AccountApiInterface.
func (AccountApi) ListWebAuthn(ctx *gin.Context) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		errorCode.GinHandler(ctx, errorCode.ErrUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, golambda.Map(account.ListWebAuthn(),
		func(_ int, cred model.Credential) api.WebAuthn {
			return api.WebAuthn{
				Id:        utils.UintToString(cred.ID),
				Name:      cred.Name,
				Os:        cred.Os,
				CreatedAt: cred.CreatedAt.Unix(),
			}
		}))
}

// CreateWebauthnTicketChallenge implements coreapi.AccountApiInterface.
func (AccountApi) CreateWebauthnTicketChallenge(ctx *gin.Context, request api.CreateTicketChallenge) {
	account := accountService.GetByEmail(request.Email)
	if account == nil {
		errorCode.GinHandler(ctx, fault.ErrUnauthorized)
		return
	}
	option, session, err := coreServices.BeginLogin(*account)
	if err != nil {
		errorCode.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":        utils.UintToString(session.ID),
		"publicKey": option.Response,
	})
}

// CreateWebAuthnChallenge implements coreapi.AccountApiInterface.
func (AccountApi) CreateWebAuthnChallenge(ctx *gin.Context) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	options, session := coreServices.BeginRegistration(*account)
	ctx.JSON(http.StatusOK, gin.H{
		"id":        utils.UintToString(session.ID),
		"publicKey": options.Response,
	})
}

type Credential struct {
	protocol.CredentialCreationResponse
	Name string `json:"name"`
	Os   string `json:"os"`
}

// CreateWebauthn implements coreapi.AccountApiInterface.
func (AccountApi) CreateWebauthn(ctx *gin.Context, id string) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	var ccr Credential
	if err := json.NewDecoder(ctx.Copy().Request.Body).Decode(&ccr); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	if err := coreServices.FinishRegistration(
		*account,
		utils.StringToUint(id),
		ccr.Name,
		ccr.Os,
		ccr.CredentialCreationResponse); err != nil {
		errorCode.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, "")
}

// GetTotp implements coreapi.AccountApiInterface.
func (AccountApi) GetTotp(ctx *gin.Context) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	ctx.JSON(http.StatusOK, api.Totp{Url: accountService.CreateTotp(account.Entity)})
}

// CreateTotp2FA implements coreapi.AccountApiInterface.
func (AccountApi) CreateTotp2FA(ctx *gin.Context, request api.PutTotpRequest) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	if err := coreServices.NFactor(account.Entity, request.Tickets, 1); err != nil {
		errorCode.GinHandler(ctx, err)
		return
	}
	url, err := accountService.UpdateTotp2FA(account.Entity, request.Url, request.TotpCode)
	if err != nil {
		errorCode.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, api.Totp{
		Url: url,
	})
}

func (AccountApi) CreateSession(c *gin.Context, createSessionRequest api.CreateSessionRequest) {
	if createSessionRequest.Tickets != nil {
		session, sessionError := accountService.CreateSessionWithTickets(*createSessionRequest.Email, *createSessionRequest.Tickets)
		if sessionError != nil {
			errorCode.GinHandler(c, sessionError)
			return
		}
		c.JSON(http.StatusCreated, api.Token{
			Secret:      session.Token,
			TokenType:   "JWT",
			TokenFormat: "bearer",
		})
	}
}

func (AccountApi) CreateAccount(ctx *gin.Context, createAccountRequest api.CreateAccountRequest) {
	if createAccountRequest.VerificationCode != nil {
		a, createError := accountService.CreateAccountWithVerificationCode(
			createAccountRequest.Email,
			*createAccountRequest.VerificationCode,
			createAccountRequest.Password)
		if createError != nil {
			errorCode.GinHandler(ctx, createError)
			return
		}
		ctx.JSON(http.StatusCreated, AccountBackward(*a))
		return
	}
	account := coreServices.CheckIfRoot(ctx)
	if account == nil {
		return
	}
	role := constants.USER
	if createAccountRequest.Role != nil {
		role = constants.Role(*createAccountRequest.Role)
	}
	a, createError := accountService.CreateAccount(
		createAccountRequest.Email,
		createAccountRequest.Password,
		role,
	)
	if createError != nil {
		errorCode.GinHandler(ctx, createError)
		return
	}
	ctx.JSON(http.StatusCreated, AccountBackward(*a))
}

func (AccountApi) ListAccount(ctx *gin.Context, ordering api.Ordering, index int64, limit int64, roles []string, email string) {

}

func (AccountApi) GetAccount(ctx *gin.Context) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	ctx.JSON(http.StatusOK, AccountBackward(*account))
}

func GetAccountById(id string) *api.Account {
	account := accountService.GetAccount(utils.StringToUint(id))
	if account == nil {
		return nil
	}
	a := AccountBackward(*account)
	return &a
}

func (a AccountApi) VerifySession(c *gin.Context, sessionVerificationRequest api.SessionVerificationRequest) {
	auth := middleware.Authorize(c)
	if auth.Status != middleware.Authorized {
		errorCode.GinHandler(c, errorCode.ErrUnauthorized)
		return
	}
	account := GetAccountById(utils.UintToString(auth.AccountId))
	if account == nil {
		errorCode.GinHandler(c, errorCode.ErrUnauthorized)
		return
	}
	c.JSON(http.StatusOK, api.Session{})
}

func (AccountApi) UpdatePassword(c *gin.Context, request api.UpdatePasswordRequest) {
	err := accountService.UpdatePassword(request.Email, request.Password, request.VerificationCode, request.NewPassword)
	if err != nil {
		errorCode.GinHandler(c, err)
		return
	}
	c.String(http.StatusNoContent, "")
}
