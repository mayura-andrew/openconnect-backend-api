package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/ideas", app.createIdeaHandler)
	router.HandlerFunc(http.MethodGet, "/v1/ideas/:id", app.showIdeaHandler)

	return router

}
