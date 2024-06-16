package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetUserExerciseMetrics(w http.ResponseWriter, r *http.Request) {
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
	user := app.contextGetUser(r)
	exerciseMetric, err := app.Models.ExerciseMetric.GetUserExerciseMetric(user.ID, date)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			env := envelope{
				"message":        "Retrieved All Exercise Metrics for user",
				"exerciseMetric": exerciseMetric,
			}
			err = app.writeJSON(w, http.StatusOK, env, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message":        "Retrieved All Exercise Metrics for user",
		"exerciseMetric": exerciseMetric,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) CreateExerciseMetric(w http.ResponseWriter, r *http.Request) {

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

	var input struct {
		Name string `json:"name"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	exerciseMetric := &models.ExerciseMetric{
		UserID: user.ID,
		Date:   date,
		Name:   input.Name,
	}
	err = app.Models.ExerciseMetric.InsertExerciseMetric(exerciseMetric)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	env := envelope{
		"message": "Successfully Created Exercise Metrics!",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) UpdateExerciseMetric(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	exerciseMetric, err := app.Models.ExerciseMetric.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Started   *string   `json:"start"`
		Ended     *string   `json:"ended"`
		Tags      *[]string `json:"tags"`
		NoOfTimes *int      `json:"no_of_times"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Started != nil {
		exerciseMetric.Started = *input.Started
	}

	if input.Ended != nil {
		exerciseMetric.Ended = *input.Ended
	}

	if input.NoOfTimes != nil {
		exerciseMetric.NoOfTimes = *input.NoOfTimes
	}

	if input.Tags != nil {
		exerciseMetric.Tags = *input.Tags
	}

	err = app.Models.ExerciseMetric.UpdateExerciseMetric(exerciseMetric, int(id))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		case errors.Is(err, models.ErrRecordAlreadyExist):
			app.recordAlreadyExistsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	env := envelope{
		"message": "Successfully Updated Exercise Metric",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteExerciseMetric(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)
	err = app.Models.ExerciseMetric.DeleteExerciseMetric(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Exercise Metric Successfully Deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
