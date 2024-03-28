package api

import (
	"errors"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetSymptoms(w http.ResponseWriter, r *http.Request) {

	symptoms, err := app.Models.Symptoms.GetSymptoms()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":  "Retrieved All Symptoms",
		"symptoms": symptoms}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetUserSymptoms(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	symptoms, err := app.Models.Symptoms.GetUserSymptoms(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":  "Retrieved All Symptoms for User",
		"symptoms": symptoms}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) InsertUserSymptoms(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Symptoms []int `json:"symptoms"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)

	err = app.Models.Symptoms.SetUserSymptoms(input.Symptoms, user.ID)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordAlreadyExist):
			app.failedValidationResponse(w, r,
				map[string]string{
					"Symptoms": "Already exists"},
			)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message": "Successfully added User Symptoms",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteUserSymptoms(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Metrics []int `json:"symptoms"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	err = app.Models.Symptoms.DeleteUserSymptoms(user.ID, input.Metrics)
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
		"message": "Deleted Symptoms for User"}

	err = app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
