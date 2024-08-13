package api

import (
	"errors"
	"net/http"

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

	tx, err := app.Models.Transaction.BeginTx()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			app.serverErrorResponse(w, r, err)
			return
		}
		err = tx.Commit()
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}()

	for i := 0; i < len(input.Metrics); i++ {
		_ = app.Models.Metrics.SetUserMetrics(tx, input.Metrics[i], user.ID)
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

	tx, err := app.Models.Transaction.BeginTx()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			app.serverErrorResponse(w, r, err)
			return
		}
		err = tx.Commit()
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}()

	for i := 0; i < len(input.Metrics); i++ {
		_ = app.Models.Metrics.DeleteUserMetrics(tx, user.ID, input.Metrics[i])
	}

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

	date, err := app.GetDate(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	exerciseBoolResult := make(chan bool)
	symsBoolResult := make(chan bool)
	sleepBoolResult := make(chan bool)
	foodBoolResult := make(chan bool)
	medicationBoolResult := make(chan bool)
	bowelBoolResult := make(chan bool)
	urineBoolResult := make(chan bool)

	defer close(exerciseBoolResult)
	defer close(symsBoolResult)
	defer close(sleepBoolResult)
	defer close(foodBoolResult)
	defer close(medicationBoolResult)
	defer close(bowelBoolResult)
	defer close(urineBoolResult)

	app.Background(func() {
		app.Models.SymsMetric.CheckUserEntry(user.ID, date, symsBoolResult)
	})
	app.Background(func() {
		app.Models.SleepMetric.CheckUserEntry(user.ID, date, sleepBoolResult)
	})
	app.Background(func() {
		app.Models.FoodMetric.CheckUserEntry(user.ID, date, foodBoolResult)
	})
	app.Background(func() {
		app.Models.ExerciseMetric.CheckUserEntry(user.ID, date, exerciseBoolResult)
	})
	app.Background(func() {
		app.Models.MedicationMetric.CheckUserEntry(user.ID, date, medicationBoolResult)
	})
	app.Background(func() {
		app.Models.BowelMetric.CheckUserEntry(user.ID, date, bowelBoolResult)
	})
	app.Background(func() {
		app.Models.UrineMetric.CheckUserEntry(user.ID, date, urineBoolResult)
	})

	symsBool := <-symsBoolResult
	sleepBool := <-sleepBoolResult
	foodBool := <-foodBoolResult
	exerciseBool := <-exerciseBoolResult
	medicationBool := <-medicationBoolResult
	urineBool := <-urineBoolResult
	bowelBool := <-bowelBoolResult

	resultMap := make(map[string]bool)
	resultMap["symptoms"] = symsBool
	resultMap["sleep"] = sleepBool
	resultMap["food"] = foodBool
	resultMap["exercise"] = exerciseBool
	resultMap["bowel"] = bowelBool
	resultMap["medication"] = medicationBool
	resultMap["urine"] = urineBool

	env := envelope{
		"message": "retrieved Tracked Metric Status for User", "metrics_status": resultMap}

	err = app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
