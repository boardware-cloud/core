package services

import (
	errorCode "github.com/boardware-cloud/common/code"
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/model/core"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	authn  *webauthn.WebAuthn
	domain string
)

func init() {
	var err error
	domain = config.GetString("boardware.web.domain")
	wconfig := &webauthn.Config{
		RPDisplayName: "Boardware Cloud",
		RPID:          domain,
		RPOrigins:     []string{"https://" + domain},
	}
	if authn, err = webauthn.New(wconfig); err != nil {
		panic(err)
	}
}

func BeginRegistration(account Account) (*protocol.CredentialCreation, core.SessionData) {
	options, session, _ := authn.BeginRegistration(account.Entity)
	sessionData := core.SessionData{
		AccountId: account.ID(),
		Data:      core.WebAuthnSessionData(*session),
	}
	DB.Save(&sessionData)
	return options, sessionData
}

func FinishRegistration(account Account, sessionId uint, name, os string, ccr protocol.CredentialCreationResponse) error {
	response, err := ccr.Parse()
	if err != nil {
		return errorCode.ErrUnauthorized
	}
	var session core.SessionData
	ctx := DB.Find(&session, sessionId)
	if ctx.RowsAffected == 0 {
		return errorCode.ErrNotFound
	}
	user, err := authn.CreateCredential(account.Entity, webauthn.SessionData(session.Data), response)
	if err != nil {
		return errorCode.ErrUnauthorized
	}
	credential := core.Credential{AccountId: account.ID(), Name: name, Os: os, Credential: core.WebAuthnCredential(*user)}
	DB.Save(&credential)
	return nil
}

func BeginLogin(account Account) (*protocol.CredentialAssertion, *core.SessionData, error) {
	options, session, err := authn.BeginLogin(account.Entity)
	if err != nil {
		return nil, nil, errorCode.ErrUnauthorized
	}
	sessionData := core.SessionData{
		AccountId: account.ID(),
		Data:      core.WebAuthnSessionData(*session),
	}
	sessionDataRepository.Save(&sessionData)
	return options, &sessionData, nil
}

func CompleteLogin(sessionId uint, car *protocol.ParsedCredentialAssertionData) (string, error) {
	session := sessionDataRepository.GetById(sessionId)
	if session == nil {
		return "", errorCode.ErrUnauthorized
	}
	account := accountRepository.GetById(session.AccountId)
	if account == nil {
		return "", errorCode.ErrUnauthorized
	}
	_, err := authn.ValidateLogin(account, webauthn.SessionData(session.Data), car)
	if err != nil {
		return "", errorCode.ErrUnauthorized
	}
	return TicketString(ticketRepository.CreateTicket("WEBAUTHN", account.ID())), nil
}
