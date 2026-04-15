package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// FirebaseAuthenticatorRESTAdaptor exposes Firebase token authentication over a REST API.
type FirebaseAuthenticatorRESTAdaptor struct {
	authenticator FirebaseAuthenticatorService
}

// NewFirebaseAuthenticatorRESTAdaptor returns a new FirebaseAuthenticatorRESTAdaptor.
func NewFirebaseAuthenticatorRESTAdaptor(
	authenticator FirebaseAuthenticatorService,
) *FirebaseAuthenticatorRESTAdaptor {
	return &FirebaseAuthenticatorRESTAdaptor{
		authenticator: authenticator,
	}
}

// FirebaseLoginRESTRequest is the expected JSON body for Firebase token login.
type FirebaseLoginRESTRequest struct {
	FirebaseToken string `json:"firebaseToken"`
}

// FirebaseLoginRESTResponse is the JSON response after a successful Firebase login.
type FirebaseLoginRESTResponse struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
}

// Login handles POST requests to authenticate a user with a Firebase ID token.
// On success, the access token is set as an HTTP-only cookie.
func (a *FirebaseAuthenticatorRESTAdaptor) Login(w http.ResponseWriter, r *http.Request) {
	var request FirebaseLoginRESTRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to decode firebase login request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	resp, err := a.authenticator.AuthenticateWithFirebaseToken(
		r.Context(),
		FirebaseAuthRequest{
			FirebaseToken: request.FirebaseToken,
		})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to authenticate with firebase token")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "authentication failed"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FirebaseLoginRESTResponse{
		UserID: resp.UserID,
		Email:  resp.Email,
	})
}
