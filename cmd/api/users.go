package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
	"github.com/olagookundavid/itoju/internal/validator"
)

func (app *Application) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Dob       time.Time `json:"dob"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := &models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Dob:       input.Dob,
		Email:     input.Email,
		Activated: true}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	if models.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.Models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	_, err = app.Models.Tokens.New(user.ID, 3*24*time.Hour, models.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// app.background(func() {
	// 	data := map[string]any{
	// 		"activationToken": token.Plaintext,
	// 		"userID":          user.ID}
	// 	err = app.mailer.Send(user.Email, "user_welcome.html", data)
	// 	if err != nil {
	// 		app.logger.PrintError(err, nil)
	// 	}
	// })
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
