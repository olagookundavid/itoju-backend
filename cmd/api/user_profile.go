package api

import (
	"errors"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	env := envelope{
		"message": "Retrieved User Profile",
		"user":    user}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateUserProfilePicHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Pic_no int `json:"pic_no"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Pic_no <= 0 {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	println(user.PicNo, input.Pic_no)
	user.PicNo = input.Pic_no
	err = app.Models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{"message": "your Profile pic as been updated"}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
