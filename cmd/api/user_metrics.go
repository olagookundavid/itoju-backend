package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) SetUserMetrics(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Metrics []int `json:"metrics"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)

	err = app.Models.Metrics.SetUserMetrics(input.Metrics, user.ID)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordAlreadyExist):
			app.failedValidationResponse(w, r,
				map[string]string{
					"Metrics": "Already exists"},
			)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message": "Successfully added track metrics",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetTrackedMetrics(w http.ResponseWriter, r *http.Request) {

	metrics, err := app.Models.Metrics.GetMetrics()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message": "Retrieved All Trackable Metrics",
		"metrics": metrics}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetUserTrackedMetrics(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	metrics, err := app.Models.Metrics.GetUserMetrics(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message": "Retrieved All Tracked Metrics for User",
		"metrics": metrics}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteUserTrackedMetrics(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Metrics []int `json:"metrics"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.Models.Metrics.DeleteUserMetrics(user.ID, input.Metrics)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message": "Deleted Tracked Metric for User"}

	err = app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetTrackedMetricsStatus(w http.ResponseWriter, r *http.Request) {

	dateString, err := app.readStringParam(r, "date")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid date format"))
		return
	}

	user := app.contextGetUser(r)
	symsBoolResult := make(chan bool)

	go func() {
		symsBool := app.Models.SymsMetric.CheckUserEntry(user.ID, date)
		symsBoolResult <- symsBool
	}()
	symsBool := <-symsBoolResult

	resultMap := make(map[string]bool)
	resultMap["symptoms"] = symsBool
	env := envelope{
		"message": "retrieved Tracked Metric Status for User", "metrics_status": resultMap}

	err = app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
