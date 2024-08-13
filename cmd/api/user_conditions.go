package api

import (
	"net/http"
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

	for i := 0; i < len(input.Conditions); i++ {
		_ = app.Models.Conditions.SetUserConditions(tx, input.Conditions[i], user.ID)
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
		Conditions []int `json:"conditions"`
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

	for i := 0; i < len(input.Conditions); i++ {
		_ = app.Models.Conditions.DeleteUserConditions(tx, user.ID, input.Conditions[i])
	}

	env := envelope{
		"message": "Deleted Conditions for User"}

	err = app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
