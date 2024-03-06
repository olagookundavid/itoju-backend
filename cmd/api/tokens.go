package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
	"github.com/olagookundavid/itoju/internal/validator"
)

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	models.ValidateEmail(v, input.Email)
	models.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.Models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	token, err := app.Models.Tokens.New(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Encode the token to JSON and send it in the response along with a 201 Created // status code.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Successfully logged in User", "data": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) CreatePasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the user's email address.
	var input struct {
		Email string `json:"email"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if models.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Try to retrieve the corresponding user record for the email address. If it can't
	// be found, return an error message to the client.
	user, err := app.Models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return an error message if the user is not activated.
	if !user.Activated {
		v.AddError("email", "user account must be activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Otherwise, create a new password reset token with a 45-minute expiry time.
	token, err := app.Models.Tokens.New(user.ID, 45*time.Minute, models.ScopePasswordReset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Email the user with their password reset token.
	app.Background(func() {
		data := map[string]any{
			"passwordResetToken": token.Plaintext,
		}

		println(data)
		// Send mail with Otp to reset mail
		// err = app.mailer.Send(user.Email, "token_password_reset.tmpl", data)
		// if err != nil {
		// 	app.logger.PrintError(err, nil)
		// }
	})
	// Send a 202 Accepted response and confirmation message to the client.
	//for now token is sent in body
	env := envelope{
		"message": "An email will be sent to you containing password reset instructions",
		"token":   token.Plaintext,
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// go func(metricID int) {
// 	defer wg.Done()
// 	defer func() {
// 		if err := recover(); err != nil {
// 			logger.PrintError(fmt.Errorf("%s", err), nil)
// 		}
// 	}()
// 	_, err := m.DB.ExecContext(ctx, query, userID, metricID)
// 	if err != nil {
// 		errors <- err
// 		return
// 	}
// 	done <- true
// 	wg.Wait()
// }(metricID)

// logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
