package routes

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/olagookundavid/itoju/cmd/api"
)

func Routes(app *api.Application) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedResponse)

	//Healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthcheckHandler)

	//Users auth
	router.HandlerFunc(http.MethodPost, "/v1/login", app.LoginHandler)
	router.HandlerFunc(http.MethodPost, "/v1/register", app.RegisterUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/password-reset", app.CreatePasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.UpdateUserPasswordHandler)
	router.Handler(http.MethodPut, "/v1/users/change-password", app.RequireActivatedAndAuthedUser(app.ChangeUserPasswordHandler))
	//Metrics
	router.Handler(http.MethodGet, "/v1/debug/vars", expvar.Handler())

	//Middleware
	//remove metric in prod!
	return app.Metrics(app.RecoverPanic(app.RateLimit(app.Authenticate(router))))
}
