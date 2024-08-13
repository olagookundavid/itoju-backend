package api

import (
	"net/http"
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

	for i := 0; i < len(input.Symptoms); i++ {
		_ = app.Models.Symptoms.SetUserSymptoms(tx, input.Symptoms[i], user.ID)
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
		Symptoms []int `json:"symptoms"`
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

	for i := 0; i < len(input.Symptoms); i++ {
		_ = app.Models.Symptoms.DeleteUserSymptoms(tx, user.ID, input.Symptoms[i])
	}
	env := envelope{
		"message": "Deleted Symptoms for User"}

	err = app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
