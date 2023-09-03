package services

import (
	"github.com/boardware-cloud/common/errors"
	"github.com/boardware-cloud/model/core"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/spf13/viper"
)

var (
	authn  *webauthn.WebAuthn
	domain string
)

func init() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	domain = viper.GetString("boardware.web.domain")
	wconfig := &webauthn.Config{
		RPDisplayName: "Boardware Cloud",                                      // Display Name for your site
		RPID:          domain,                                                 // Generally the FQDN for your site
		RPOrigins:     []string{"https://" + domain, "http://localhost:3000"}, // The origin URLs allowed for WebAuthn requests
	}
	if authn, err = webauthn.New(wconfig); err != nil {
		panic(err)
	}
}

func DeleteWebAuthn(account core.Account, id uint) *errors.Error {
	ctx := DB.Where("account_id = ? AND id = ?", account.ID, id).Delete(&core.Credential{})
	if ctx.RowsAffected == 0 {
		return errors.NotFoundError()
	}
	return nil
}

func ListWebAuthn(account core.Account) []core.Credential {
	var webauthns []core.Credential = make([]core.Credential, 0)
	DB.Where("account_id = ?", account.ID).Find(&webauthns)
	return webauthns
}

func BeginRegistration(account core.Account) (*protocol.CredentialCreation, core.SessionData) {
	options, session, _ := authn.BeginRegistration(account)
	sessionData := core.SessionData{
		AccountId: account.ID,
		Data:      core.WebAuthnSessionData(*session),
	}
	DB.Save(&sessionData)
	return options, sessionData
}

func FinishRegistration(account core.Account, sessionId uint, name, os string, ccr protocol.CredentialCreationResponse) *errors.Error {
	response, err := ccr.Parse()
	if err != nil {
		return errors.AuthenticationError()
	}
	var session core.SessionData
	ctx := DB.Find(&session, sessionId)
	if ctx.RowsAffected == 0 {
		return errors.NotFoundError()
	}
	user, err := authn.CreateCredential(account, webauthn.SessionData(session.Data), response)
	if err != nil {
		return errors.AuthenticationError()
	}
	credential := core.Credential{AccountId: account.ID, Name: name, Os: os, Credential: core.WebAuthnCredential(*user)}
	DB.Save(&credential)
	return nil
}

func BeginLogin(account core.Account) (*protocol.CredentialAssertion, *core.SessionData, *errors.Error) {
	options, session, err := authn.BeginLogin(account)
	if err != nil {
		return nil, nil, errors.NotFoundError()
	}
	sessionData := core.SessionData{
		AccountId: account.ID,
		Data:      core.WebAuthnSessionData(*session),
	}
	DB.Save(&sessionData)
	return options, &sessionData, nil
}

func FinishLogin(sessionId uint, car *protocol.ParsedCredentialAssertionData) (string, *errors.Error) {
	var session core.SessionData
	ctx := DB.Find(&session, sessionId)
	if ctx.RowsAffected == 0 {
		return "", errors.NotFoundError()
	}
	account, errg := core.GetAccount(session.AccountId)
	if errg != nil {
		return "", errg
	}
	_, err := authn.ValidateLogin(account, webauthn.SessionData(session.Data), car)
	if err != nil {
		return "", errors.AuthenticationError()
	}
	return TicketString(createTicket("WEBAUTHN", account.ID)), nil
}
