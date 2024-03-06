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

	//Profile
	router.Handler(http.MethodGet, "/v1/users/profile", app.RequireActivatedAndAuthedUser(app.GetUserProfileHandler))

	//User tracked metrics
	router.Handler(http.MethodPost, "/v1/user/metrics", app.RequireActivatedAndAuthedUser(app.SetUserMetrics))
	router.HandlerFunc(http.MethodGet, "/v1/allmetrics", (app.GetTrackedMetrics))
	router.Handler(http.MethodGet, "/v1/user/metrics", app.RequireActivatedAndAuthedUser((app.GetUserTrackedMetrics)))
	router.Handler(http.MethodDelete, "/v1/user/metrics", app.RequireActivatedAndAuthedUser((app.DeleteUserTrackedMetrics)))

	//User smileys
	router.HandlerFunc(http.MethodGet, "/v1/allsmileys", (app.GetSmileys))
	router.Handler(http.MethodGet, "/v1/user/smileys", app.RequireActivatedAndAuthedUser((app.GetUserSmileys)))
	router.Handler(http.MethodPost, "/v1/user/smileys", app.RequireActivatedAndAuthedUser((app.InsertUserSmileys)))
	router.Handler(http.MethodGet, "/v1/user/smileys_count/:id", app.RequireActivatedAndAuthedUser((app.GetUserSmileysCount)))

	//Metrics
	router.Handler(http.MethodGet, "/v1/debug/vars", expvar.Handler())

	//Middleware
	return app.Metrics(app.RecoverPanic(app.RateLimit(app.Authenticate(router))))
}
