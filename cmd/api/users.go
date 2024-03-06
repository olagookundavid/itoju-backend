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
	err = app.writeJSON(w, http.StatusCreated, envelope{
		"message": "Successful Registered User",
		"user":    user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the user's new password and password reset token.
	var input struct {
		Password       string `json:"password"`
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	models.ValidatePasswordPlaintext(v, input.Password)
	models.ValidateTokenPlaintext(v, input.TokenPlaintext)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.Models.Users.GetForToken(models.ScopePasswordReset, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			v.AddError("token", "invalid or expired password reset token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Set the new password for the user.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Save the updated user record in our database, checking for any edit conflicts as // normal.
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
	// If everything was successful, then delete all password reset tokens for the user.
	err = app.Models.Tokens.DeleteAllForUser(models.ScopePasswordReset, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the user a confirmation message.
	env := envelope{"message": "your password was successfully reset"}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) ChangeUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the user's new password and password reset token.
	var input struct {
		Password           string `json:"password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.NewPassword == input.ConfirmNewPassword, "new password", "doesn't match confirm password")
	v.Check(input.NewPassword != input.Password, "password", "old and new password are the same")
	models.ValidatePasswordPlaintext(v, input.Password)
	models.ValidatePasswordPlaintext(v, input.ConfirmNewPassword)
	models.ValidatePasswordPlaintext(v, input.NewPassword)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user := app.contextGetUser(r)
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	err = user.Password.Set(input.NewPassword)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Save the updated user record in our database, checking for any edit conflicts as normal
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
	// Send the user a confirmation message.
	env := envelope{"message": "your password was successfully changed"}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	print("successful")

	env := envelope{
		"message": "Retrieved User Profile",
		"user":    user}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
