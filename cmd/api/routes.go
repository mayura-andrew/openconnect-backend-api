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

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}
