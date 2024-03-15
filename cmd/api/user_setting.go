package api

import (
	"errors"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
	"github.com/olagookundavid/itoju/internal/validator"
)

func (app *Application) InsertMenses(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Period_len int `json:"period_len"`
		Cycle_len  int `json:"cycle_len"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	menses := &models.Menses{
		Id: user.ID, Period_len: input.Period_len, Cycle_len: input.Cycle_len,
	}
	v := validator.New()
	if models.ValidateMenses(v, menses); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Models.Menses.InsertMenses(menses)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	env := envelope{
		"message": "Successfully added Menstruation Data ",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) GetMenses(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	menses, err := app.Models.Menses.GetMenses(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotYetSet(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message": "Retrieved User Menses",
		"menses":  menses}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateMenses(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	menses, err := app.Models.Menses.GetMenses(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotYetSet(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Period_len *int `json:"period_len"`
		Cycle_len  *int `json:"cycle_len"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Period_len != nil {
		menses.Period_len = *input.Period_len
	}
	if input.Cycle_len != nil {
		menses.Cycle_len = *input.Cycle_len
	}

	v := validator.New()
	if models.ValidateMenses(v, menses); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.Models.Menses.UpdateMenses(menses)
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
		"message": "Successfully updated Menstruation",
		"menses":  menses,
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//Body Measure

func (app *Application) InsertBodyMeasure(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Height int `json:"height"`
		Weight int `json:"weight"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	bodyMeasure := &models.BodyMeasure{
		Id: user.ID, Height: input.Height, Weight: input.Weight,
	}
	v := validator.New()
	if models.ValidateBodyMeasure(v, bodyMeasure); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Models.BodyMeasure.InsertBodyMeasure(bodyMeasure)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	env := envelope{
		"message": "Successfully added Body Measure Data ",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) GetBodyMeasure(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	bodyMeasure, err := app.Models.BodyMeasure.GetBodyMeasure(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotYetSet(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message":      "Retrieved User Body Measure",
		"body_measure": bodyMeasure}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateBodyMeasure(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	bodyMeasure, err := app.Models.BodyMeasure.GetBodyMeasure(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotYetSet(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Height *int `json:"height"`
		Weight *int `json:"weight"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Height != nil {
		bodyMeasure.Height = *input.Height
	}
	if input.Weight != nil {
		bodyMeasure.Weight = *input.Weight
	}

	v := validator.New()
	if models.ValidateBodyMeasure(v, bodyMeasure); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.Models.BodyMeasure.UpdateBodyMeasure(bodyMeasure)
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
		"message":      "Successfully updated Body Measure",
		"body_measure": bodyMeasure,
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
