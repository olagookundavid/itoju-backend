package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetMenstrualCycle(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	id, err := app.Models.UserPeriod.GetMensesCycleId(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	periodDays, err := app.Models.UserPeriod.GetCycleDays(id, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":     "Retrieved All Period data",
		"period_days": periodDays}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetCycleDay(w http.ResponseWriter, r *http.Request) {
	id, err := app.readStringParam(r, "id")
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}
	periodDay, err := app.Models.UserPeriod.GetCycleDay(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":    "Retrieved Period data",
		"period_day": periodDay}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) AddMenstrualCycle(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	var input struct {
		StartDate    string `json:"start_date"`
		CycleLength  int    `json:"cycle_length"`
		PeriodLength int    `json:"period_length"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	date, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		err := errors.New("invalid date format")
		app.badRequestResponse(w, r, err)
		return
	}

	// Start a transaction
	tx, err := app.Models.Transaction.BeginTx()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}()

	// Insert the menstrual cycle within the transaction
	cycle := app.Models.UserPeriod.ReturnMenstrualCycle(
		user.ID,
		input.CycleLength,
		input.PeriodLength,
		date,
	)

	cycleID, err := app.Models.UserPeriod.InsertMenstrualCycleTx(tx, &cycle)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordAlreadyExist):
			app.recordAlreadyExistsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Insert cycle days
	for i := 0; i < input.PeriodLength; i++ {
		day := app.Models.UserPeriod.ReturnCycleDay(cycleID, user.ID, true, false, cycle.StartDate.AddDate(0, 0, i))
		if err := app.Models.UserPeriod.InsertCycleDayTx(tx, &day); err != nil {
			return
		}
	}
	// Insert regular days
	// for i := input.PeriodLength; i < (input.PeriodLength + 9); i++ {
	// 	day := app.Models.UserPeriod.ReturnCycleDay(cycleID, user.ID, false, false, cycle.StartDate.AddDate(0, 0, i))
	// 	if err := app.Models.UserPeriod.InsertCycleDayTx(tx, &day); err != nil {
	// 		return
	// 	}
	// }
	// Insert ovulation days
	for i := (input.PeriodLength + 9); i < (input.CycleLength); i++ {
		day := app.Models.UserPeriod.ReturnCycleDay(cycleID, user.ID, false, true, cycle.StartDate.AddDate(0, 0, i))
		if err := app.Models.UserPeriod.InsertCycleDayTx(tx, &day); err != nil {
			return
		}
	}

	env := envelope{
		"message": "Successful Created User Cycle",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateMenstrualCycle(w http.ResponseWriter, r *http.Request) {
	id, err := app.readStringParam(r, "id")
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	cycleDay, err := app.Models.UserPeriod.GetCycleDay(id)
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
		IsPeriod    *bool     `json:"is_period"`
		IsOvulation *bool     `json:"is_ovulation"`
		Flow        *float32  `json:"flow"`
		Pain        *float32  `json:"pain"`
		Tags        *[]string `json:"tags"`
		CMQ         *string   `json:"cmq"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Pain != nil {
		cycleDay.Pain = *input.Pain
	}
	if input.Flow != nil {
		cycleDay.Flow = *input.Flow
	}
	if input.IsOvulation != nil {
		cycleDay.IsOvulation = *input.IsOvulation
	}
	if input.IsPeriod != nil {
		cycleDay.IsPeriod = *input.IsPeriod
	}
	if input.CMQ != nil {
		cycleDay.CMQ = *input.CMQ
	}
	if input.Tags != nil {
		cycleDay.Tags = *input.Tags
	}

	err = app.Models.UserPeriod.UpdateCycleDay(cycleDay)
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
		"message":  "Successfully updated Cycle Day",
		"cycleDay": cycleDay,
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
