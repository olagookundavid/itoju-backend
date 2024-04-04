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
	router.Handler(http.MethodGet, "/v1/user/metrics_status/:date", app.RequireActivatedAndAuthedUser((app.GetTrackedMetricsStatus)))

	//User smileys
	router.HandlerFunc(http.MethodGet, "/v1/allsmileys", (app.GetSmileys))
	router.Handler(http.MethodGet, "/v1/user/smileys", app.RequireActivatedAndAuthedUser((app.GetUserSmileys)))
	router.Handler(http.MethodGet, "/v1/user/lastestsmileys", app.RequireActivatedAndAuthedUser((app.GetLatestUserSmileyForToday)))
	router.Handler(http.MethodPost, "/v1/user/smileys", app.RequireActivatedAndAuthedUser((app.InsertUserSmileys)))
	router.Handler(http.MethodGet, "/v1/user/smileys_count/:id", app.RequireActivatedAndAuthedUser((app.GetUserSmileysCount)))

	//User symptoms
	router.HandlerFunc(http.MethodGet, "/v1/allsymptoms", (app.GetSymptoms))
	// router.Handler(http.MethodGet, "/v1/user/symptoms", app.RequireActivatedAndAuthedUser((app.GetUserSymptoms)))
	// router.Handler(http.MethodPost, "/v1/user/symptoms", app.RequireActivatedAndAuthedUser((app.InsertUserSymptoms)))
	// router.Handler(http.MethodDelete, "/v1/user/symptoms", app.RequireActivatedAndAuthedUser((app.DeleteUserSymptoms)))

	//User conditions
	router.HandlerFunc(http.MethodGet, "/v1/allconditions", (app.GetConditions))
	router.Handler(http.MethodGet, "/v1/user/conditions", app.RequireActivatedAndAuthedUser((app.GetUserConditions)))
	router.Handler(http.MethodPost, "/v1/user/conditions", app.RequireActivatedAndAuthedUser((app.InsertUserConditions)))
	router.Handler(http.MethodDelete, "/v1/user/conditions", app.RequireActivatedAndAuthedUser((app.DeleteUserConditions)))

	//Resources
	router.HandlerFunc(http.MethodGet, "/v1/resources", (app.GetResources))
	router.HandlerFunc(http.MethodPost, "/v1/resources", (app.InsertResources))
	router.HandlerFunc(http.MethodPut, "/v1/resources/:id", (app.UpdateResources))
	router.HandlerFunc(http.MethodDelete, "/v1/resources/:id", (app.DeleteResources))

	//Setting
	// router.Handler(http.MethodPost, "/v1/user/menses", app.RequireActivatedAndAuthedUser((app.InsertMenses)))
	router.Handler(http.MethodGet, "/v1/user/menses", app.RequireActivatedAndAuthedUser((app.GetMenses)))
	router.Handler(http.MethodPut, "/v1/user/menses", app.RequireActivatedAndAuthedUser((app.UpdateMenses)))
	// router.Handler(http.MethodPost, "/v1/user/bodymeasure", app.RequireActivatedAndAuthedUser((app.InsertBodyMeasure)))
	router.Handler(http.MethodGet, "/v1/user/bodymeasure", app.RequireActivatedAndAuthedUser((app.GetBodyMeasure)))
	router.Handler(http.MethodPut, "/v1/user/bodymeasure", app.RequireActivatedAndAuthedUser((app.UpdateBodyMeasure)))

	//SymsMetric
	router.Handler(http.MethodPost, "/v1/user/symsMetric", app.RequireActivatedAndAuthedUser((app.CreateSymsMetric)))
	router.Handler(http.MethodPut, "/v1/user/symsMetric/:id", app.RequireActivatedAndAuthedUser((app.UpdateSymsMetric)))
	router.Handler(http.MethodDelete, "/v1/user/symsMetric/:id", app.RequireActivatedAndAuthedUser((app.DeleteSymsMetric)))
	router.Handler(http.MethodGet, "/v1/user/symsMetric/:date", app.RequireActivatedAndAuthedUser((app.GetUserSymsMetric)))
	router.Handler(http.MethodGet, "/v1/user/symsN/:id", app.RequireActivatedAndAuthedUser((app.GetUserTopNSyms)))

	//Achievement
	router.Handler(http.MethodGet, "/v1/user/getDaysTracked", app.RequireActivatedAndAuthedUser((app.GetDaysTrackedInARow)))
	router.Handler(http.MethodGet, "/v1/user/getDaysTrackedFree", app.RequireActivatedAndAuthedUser((app.GetDaysTrackedFree)))

	//Metrics
	router.Handler(http.MethodGet, "/v1/debug/vars", expvar.Handler())

	//Middleware
	return app.Metrics(app.RecoverPanic(app.RateLimit(app.Authenticate(router))))
}
