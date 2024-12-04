package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/OpenConnectOUSL/backend-api-v1/internal/data"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID string `json:"id"`
	Email string `json:"email"`
	VerifiedEmail bool `json:"verified_email"`
	Name string `json:"name"`
	Picture string `json:"picture"`
}

func (app *application) initGoogleOAuth() {
    app.googleOauthConfig = &oauth2.Config{
        RedirectURL:  app.config.oauth.redirectURL,
        ClientID:     app.config.oauth.googleClientID,
        ClientSecret: app.config.oauth.googleClientSecret,
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
        Endpoint: google.Endpoint,
    }
}

func (app *application) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateToken()

	// Store state token in cookie
	http.SetCookie(w, &http.Cookie{
		Name: "oauthState",
	Value: state,
	Expires: time.Now().Add(15 * time.Minute),
	HttpOnly: true,
	Secure: true,
	SameSite: http.SameSiteStrictMode,
})

	url := app.googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	// Verify state token
	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != state {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// Exchange code for token
	token, err := app.googleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get user info from Google
	googleUser, err := app.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//Find or create user
	user, err := app.models.Users.FindOrCreateFromGoogle(googleUser)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Generate authentication token
	authToken, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": authToken}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func generateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (app *application) getGoogleUserInfo(accessToken string) (*data.GoogleUser, error) {
    resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err 
	}

	defer resp.Body.Close()

	var user data.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
