package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) CreateSymsMetric(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id int `json:"symsptom_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	symsMetric := &models.SymsMetric{
		Id: input.Id,
	}

	err = app.Models.SymsMetric.CreateSymsMetric(user.ID, *symsMetric)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordAlreadyExist):
			app.recordAlreadyExistsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	env := envelope{
		"message": "Successfully added Symptom",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) GetUserSymsMetric(w http.ResponseWriter, r *http.Request) {

	dateString, err := app.readStringParam(r, "date")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Parse the date string into a time.Time object
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid date format"))
		return
	}
	user := app.contextGetUser(r)
	symsMetric, err := app.Models.SymsMetric.GetUserSymptomsMetric(user.ID, date)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":    "Retrieved All Symptom Metrics for user",
		"symsMetric": symsMetric}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateSymsMetric(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	symsMetric, err := app.Models.SymsMetric.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		MorningSeverity   *float32 `json:"morning_severity"`
		AfternoonSeverity *float32 `json:"afternoon_severity"`
		NightSeverity     *float32 `json:"night_severity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.MorningSeverity != nil {
		symsMetric.MorningSeverity = *input.MorningSeverity
	}
	if input.AfternoonSeverity != nil {
		symsMetric.AfternoonSeverity = *input.AfternoonSeverity
	}
	if input.NightSeverity != nil {
		symsMetric.NightSeverity = *input.NightSeverity
	}

	err = app.Models.SymsMetric.UpdateSymsMetric(symsMetric, int(id))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		case errors.Is(err, models.ErrRecordAlreadyExist):
			app.recordAlreadyExistsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	env := envelope{
		"message": "Successfully updated Symptom Metric",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteSymsMetric(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)
	err = app.Models.SymsMetric.DeleteSymsMetric(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Symptom Metric successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
