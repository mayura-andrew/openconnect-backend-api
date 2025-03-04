package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/ideas", app.requirePermission("ideas:read", app.listIdeasHandler))
	router.HandlerFunc(http.MethodPost, "/v1/ideas", app.requirePermission("ideas:write", app.createIdeaHandler))
	router.HandlerFunc(http.MethodGet, "/v1/ideas/:id", app.requirePermission("ideas:read", app.showIdeaHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/ideas/:id", app.requirePermission("ideas:write", app.updateIdeaHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/ideas/:id", app.requirePermission("ideas:write", app.deleteIdeaHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password-reset", app.updateUserPasswordHandler)

	router.HandlerFunc(http.MethodPost, "/v1/auth/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/tokens/password-reset-request", app.createPasswordResetTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/auth/google/login", app.googleLoginHandler)
	router.HandlerFunc(http.MethodGet, "/v1/auth/google/callback", app.googleCallbackHandler)

	router.HandlerFunc(http.MethodGet, "/v1/profile", app.requirePermission("ideas:read", app.getProfileHandler))
	router.HandlerFunc(http.MethodGet, "/v1/profiles/search", app.requirePermission("ideas:read", app.searchProfilesHandler))
	router.HandlerFunc(http.MethodPut, "/v1/profile/new", app.requirePermission("ideas:write", app.createProfileHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/profile/update", app.requirePermission("ideas:write", app.updateProfileHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/profile/delete", app.requirePermission("ideas:write", app.deleteProfileHandler))
	//router.HandlerFunc(http.MethodGet, "/v1/profiles/:username", app.requirePermission("ideas:read", app.getProfileByUsernameHandler))
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}
