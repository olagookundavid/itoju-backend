package routes

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/olagookundavid/itoju/cmd/api"
)

func Routes(app *api.Application) http.Handler {
	router := httprouter.New()
	// router.NotFound = http.HandlerFunc(app.notFoundResponse)
	// router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//Healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthcheckHandler)

	//Users auth
	router.HandlerFunc(http.MethodPost, "/v1/login", app.LoginHandler)
	router.HandlerFunc(http.MethodPost, "/v1/register", app.RegisterUserHandler)
	//metrics
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	//
	return router
}
