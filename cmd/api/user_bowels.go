package api

import (
	"errors"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetUserBowelMetrics(w http.ResponseWriter, r *http.Request) {

	date, err := app.GetDate(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	bowelMetric, err := app.Models.BowelMetric.GetUserBowelMetrics(user.ID, date)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":      "Retrieved All Bowel Metrics for user",
		"bowelMetrics": bowelMetric}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) UpdateBowelMetric(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}
	bowelMetric, err := app.Models.BowelMetric.GetUserBowelMetric(user.ID, id)
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
		Time *string   `json:"time"`
		Type *float64  `json:"type"`
		Pain *float64  `json:"pain"`
		Tags *[]string `json:"tags"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Time != nil {
		bowelMetric.Time = *input.Time
	}
	if input.Type != nil {
		bowelMetric.Type = *input.Type
	}
	if input.Pain != nil {
		bowelMetric.Pain = *input.Pain
	}
	if input.Tags != nil {
		bowelMetric.Tags = *input.Tags
	}

	err = app.Models.BowelMetric.UpdateBowelMetric(bowelMetric)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	env := envelope{
		"message": "Successfully updated User Bowel Metrics",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) CreateBowelMetric(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	date, err := app.GetDate(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Time string   `json:"time"`
		Type float64  `json:"type"`
		Pain float64  `json:"pain"`
		Tags []string `json:"tags"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	bowelMetric := &models.BowelMetric{
		Time: input.Time, Type: input.Type, Pain: input.Pain, Tags: input.Tags, Date: date}

	err = app.Models.BowelMetric.InsertBowelMetric(user.ID, bowelMetric)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	env := envelope{
		"message": "Successfully Created User Bowel Metrics!",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteBowelMetric(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	err = app.Models.BowelMetric.DeleteBowelMetric(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Bowel Metric successfully deleted"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
