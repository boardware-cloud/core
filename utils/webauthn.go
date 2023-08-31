package utils

import (
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	w   *webauthn.WebAuthn
	err error
)

func init() {
	wconfig := &webauthn.Config{
		RPDisplayName: "Boardware Cloud",                       // Display Name for your site
		RPID:          "cloud.boardware.com",                   // Generally the FQDN for your site
		RPOrigins:     []string{"https://cloud.boardware.com"}, // The origin URLs allowed for WebAuthn requests
	}

	if w, err = webauthn.New(wconfig); err != nil {
		panic(err)
	}
	// protocol.ParseCredentialCreationResponseBody
	// w.BeginRegistration(user)
	// w.BeginRegistration(user)

}

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	// user := datastore.GetUser() // Find or create the new user
	// options, session, err := w.BeginRegistration(user)
	// // handle errors if present
	// // store the sessionData values
	// JSONResponse(w, options, http.StatusOK) // return the options generated
	// options.publicKey contain our registration options
}
