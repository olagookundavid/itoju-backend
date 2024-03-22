package api

import (
	"errors"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetConditions(w http.ResponseWriter, r *http.Request) {

	conditions, err := app.Models.Conditions.GetConditions()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":    "Retrieved All Conditions",
		"conditions": conditions}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetUserConditions(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	conditions, err := app.Models.Conditions.GetUserConditions(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":    "Retrieved All Conditions for User",
		"conditions": conditions}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) InsertUserConditions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Conditions []int `json:"conditions"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	err = app.Models.Conditions.SetUserConditions(input.Conditions, user.ID)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordAlreadyExist):
			app.failedValidationResponse(w, r,
				map[string]string{
					"Conditions": "Already exists"},
			)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	env := envelope{
		"message": "Successfully added User Conditions",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteUserConditions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Metrics []int `json:"conditions"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.Models.Conditions.DeleteUserConditions(user.ID, input.Metrics)
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
		"message": "Deleted Conditions for User"}

	err = app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
