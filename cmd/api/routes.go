package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandlerFunc(http.MethodDelete, "/v1/module-infos/:id", app.deleteModuleInfo)
	router.HandlerFunc(http.MethodPut, "/v1/module-infos/:id", app.editModuleInfo)
	router.HandlerFunc(http.MethodGet, "/v1/module-infos", app.getLastFiftyModuleInfo)
	router.HandlerFunc(http.MethodGet, "/v1/module-infos/:id", app.getModuleInfo)
	router.HandlerFunc(http.MethodPost, "/v1/module-infos/create", app.createModuleInfo)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/department-info", app.createDepInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/department-info/:id", app.getDepInfoHandler)
	//router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	//router.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.updateMovieHandler)
	//router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	// Return the httprouter instance.

	router.Handler(http.MethodPost, "/v1/users", app.requireAdminRole(app.registerUserInfoHandler))
	router.Handler(http.MethodPut, "/v1/users/activated", app.requireAdminRole(app.activateUserHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.Handler(http.MethodPut, "/v1/users/edit", app.requireAdminRole(app.editUserInfo))
	router.Handler(http.MethodDelete, "/v1/users/delete", app.requireAdminRole(app.deleteUserInfo))
	router.Handler(http.MethodGet, "/v1/users/:id", app.requireAdminRole(app.getUserInfoHandler))
	router.Handler(http.MethodGet, "/v1/users/all", app.requireAdminRole(app.listUsersHandler))

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
