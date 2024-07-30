package api

import (
	"errors"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetUserUrineMetrics(w http.ResponseWriter, r *http.Request) {

	date, err := app.GetDate(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	urineMetric, err := app.Models.UrineMetric.GetUserUrineMetrics(user.ID, date)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":      "Retrieved All Urine Metrics for user",
		"urineMetrics": urineMetric}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) UpdateUrineMetric(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}
	urineMetric, err := app.Models.UrineMetric.GetUserUrineMetric(user.ID, id)
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
		Time     *string   `json:"time"`
		Type     *float64  `json:"type"`
		Pain     *float64  `json:"pain"`
		Quantity *float64  `json:"quantity"`
		Tags     *[]string `json:"tags"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Time != nil {
		urineMetric.Time = *input.Time
	}
	if input.Type != nil {
		urineMetric.Type = *input.Type
	}
	if input.Pain != nil {
		urineMetric.Pain = *input.Pain
	}
	if input.Tags != nil {
		urineMetric.Tags = *input.Tags
	}
	if input.Quantity != nil {
		urineMetric.Quantity = *input.Quantity
	}

	err = app.Models.UrineMetric.UpdateUrineMetric(urineMetric)
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
		"message": "Successfully updated User Urine Metrics",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) CreateUrineMetric(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	date, err := app.GetDate(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Time     string   `json:"time"`
		Type     float64  `json:"type"`
		Pain     float64  `json:"pain"`
		Quantity float64  `json:"quantity"`
		Tags     []string `json:"tags"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	urineMetric := &models.UrineMetric{
		Time: input.Time, Type: input.Type, Pain: input.Pain, Tags: input.Tags, Quantity: input.Quantity, Date: date}

	err = app.Models.UrineMetric.InsertUrineMetric(user.ID, urineMetric)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.Background(func() {
		_ = app.Models.UserPoint.InsertPoint(user.ID, "Urine", 2)
	})
	env := envelope{
		"message": "Successfully Created User Urine Metrics!",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteUrineMetric(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	err = app.Models.UrineMetric.DeleteUrineMetric(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Urine Metric successfully deleted"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
