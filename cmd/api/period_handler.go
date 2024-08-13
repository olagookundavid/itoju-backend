package api

import (
	"errors"
	"net/http"
	"time"
)

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
			app.serverErrorResponse(w, r, err)
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
		app.serverErrorResponse(w, r, err)
		return
	}

	// Example: Adding cycle days based on the cycle within the same transaction
	for i := 0; i < input.PeriodLength; i++ {
		day := app.Models.UserPeriod.ReturnCycleDay(cycleID, false, false, cycle.StartDate.AddDate(0, 0, i))
		// day := models.CycleDay{
		// 	CycleID:     cycleID,
		// 	Date:        cycle.StartDate.AddDate(0, 0, i),
		// 	IsPeriod:    true,
		// 	IsOvulation: false,
		// 	Flow:        3, // Example flow level
		// 	Pain:        2, // Example pain level
		// 	Tags:        []string{"cramps"},
		// 	CMQ:         "Sticky",
		// }
		_ = app.Models.UserPeriod.InsertCycleDayTx(tx, &day)

	}
	env := envelope{
		"message": "Success",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
