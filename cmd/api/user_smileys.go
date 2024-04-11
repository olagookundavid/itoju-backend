package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetSmileys(w http.ResponseWriter, r *http.Request) {

	smileys, err := app.Models.Smileys.GetSmileys()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message": "Retrieved All Smileys",
		"smileys": smileys}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) InsertUserSmileys(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SmileyID int      `json:"smiley_id"`
		Tags     []string `json:"tags"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	smiley := &models.Smileys{
		Id: input.SmileyID, Tags: input.Tags,
	}
	user := app.contextGetUser(r)
	err = app.Models.Smileys.InsertUserSmileys(user.ID, *smiley)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	env := envelope{
		"message": "Successfully added User smiley",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) GetUserSmileys(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	smileys, err := app.Models.Smileys.GetUserSmileys(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message": "Retrieved All Smileys for User",
		"smileys": smileys}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetLatestUserSmileyForToday(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
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
	smiley, err := app.Models.Smileys.GetLatestUserSmileyForToday(user.ID, date)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message": "Retrieved All Smileys for User",
		"smileys": smiley}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetUserSmileysCountInXDays(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	smileys, totalCount, err := app.Models.Smileys.GetUserSmileysCount(user.ID, int(id))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":     fmt.Sprintf("Retrieved All Smiley's count for User in %d day(s)", id),
		"smileys":     smileys,
		"total_count": totalCount,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
