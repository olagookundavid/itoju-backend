package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetUserSleepMetrics(w http.ResponseWriter, r *http.Request) {
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

	dayMetricChan := make(chan *models.SleepMetric)
	nightMetricChan := make(chan *models.SleepMetric)
	dayErrChan := make(chan error)
	nightErrChan := make(chan error)
	var daySleepMetric, nightSleepMetric *models.SleepMetric

	app.Background(func() {
		app.getUserSleepMetricAsync(user.ID, date, dayMetricChan, dayErrChan, false)
	})
	app.Background(func() {
		app.getUserSleepMetricAsync(user.ID, date, nightMetricChan, nightErrChan, true)
	})

	for i := 0; i < 2; i++ {
		select {
		case dayMetric := <-dayMetricChan:
			daySleepMetric = dayMetric
		case nightMetric := <-nightMetricChan:
			nightSleepMetric = nightMetric
		case <-dayErrChan:

		case <-nightErrChan:
		}
	}

	defer close(dayMetricChan)
	defer close(dayErrChan)
	defer close(nightMetricChan)
	defer close(nightErrChan)
	env := envelope{
		"message":          "Retrieved All Sleep Metrics for user",
		"daySleepMetric":   daySleepMetric,
		"nightSleepMetric": nightSleepMetric,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateSleepMetric(w http.ResponseWriter, r *http.Request) {
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
	isNight, err := app.readBoolParam(r, "is_night")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	sleepMetric, err := app.Models.SleepMetric.GetUserSleepMetric(user.ID, date, isNight)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			var input struct {
				TimeSlept  string   `json:"time_slept"`
				TimeWokeUp string   `json:"time_woke_up"`
				Severity   float64  `json:"severity"`
				Tags       []string `json:"tags"`
			}
			err := app.readJSON(w, r, &input)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}
			sleepMetric := &models.SleepMetric{
				IsNight: isNight, TimeSlept: input.TimeSlept, TimeWokeUp: input.TimeWokeUp, Tags: input.Tags, Date: date, Severity: input.Severity}

			err = app.Models.SleepMetric.InsertSleepMetric(user.ID, sleepMetric)

			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			env := envelope{
				"message": "Successfully updated User Sleep Metrics!",
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

	var input struct {
		TimeSlept  *string   `json:"time_slept"`
		TimeWokeUp *string   `json:"time_woke_up"`
		Severity   *float64  `json:"severity"`
		Tags       *[]string `json:"tags"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Severity != nil {
		sleepMetric.Severity = *input.Severity
	}
	if input.TimeSlept != nil {
		sleepMetric.TimeSlept = *input.TimeSlept
	}
	if input.TimeWokeUp != nil {
		sleepMetric.TimeWokeUp = *input.TimeWokeUp
	}
	if input.Tags != nil {
		sleepMetric.Tags = *input.Tags
	}

	err = app.Models.SleepMetric.UpdateSleepMetric(sleepMetric)
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
		"message": "Successfully updated User Sleep Metrics",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getUserSleepMetricAsync(userID string, date time.Time, sendMetric chan<- (*models.SleepMetric), errChan chan<- error, isNight bool) {
	metric, err := app.Models.SleepMetric.GetUserSleepMetric(userID, date, isNight)
	if err != nil {
		errChan <- err
		return
	}
	sendMetric <- metric
}
