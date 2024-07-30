package api

import (
	"errors"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetUserMedicationMetrics(w http.ResponseWriter, r *http.Request) {

	date, err := app.GetDate(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	medicationMetric, err := app.Models.MedicationMetric.GetUserMedicationMetrics(user.ID, date)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":           "Retrieved All Medication Metrics for user",
		"medicationMetrics": medicationMetric}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) UpdateMedicationMetric(w http.ResponseWriter, r *http.Request) {
	// user := app.contextGetUser(r)
	// id, err := app.readIDParam(r)
	// if err != nil {
	// 	app.NotFoundResponse(w, r)
	// 	return
	// }
	// medicationMetric, err := app.Models.MedicationMetric.GetUserMedicationMetric(user.ID, id)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, models.ErrRecordNotFound):
	// 		app.NotFoundResponse(w, r)
	// 	default:
	// 		app.serverErrorResponse(w, r, err)
	// 	}
	// 	return
	// }

	var medicationMetric *models.MedicationMetric = &models.MedicationMetric{}
	var input struct {
		Time     *string  `json:"time"`
		Name     *string  `json:"name"`
		Metric   *string  `json:"metric"`
		Dosage   *float64 `json:"dosage"`
		Quantity *float64 `json:"quantity"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Time != nil {
		medicationMetric.Time = *input.Time
	}
	if input.Dosage != nil {
		medicationMetric.Dosage = *input.Dosage
	}
	if input.Quantity != nil {
		medicationMetric.Quantity = *input.Quantity
	}
	if input.Name != nil {
		medicationMetric.Name = *input.Name
	}
	if input.Metric != nil {
		medicationMetric.Metric = *input.Metric
	}

	err = app.Models.MedicationMetric.UpdateMedicationMetric(medicationMetric)
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
		"message": "Successfully updated User Medication Metrics",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) CreateMedicationMetric(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	date, err := app.GetDate(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Time     string  `json:"time"`
		Name     string  `json:"name"`
		Metric   string  `json:"metric"`
		Dosage   float64 `json:"dosage"`
		Quantity float64 `json:"quantity"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	medicationMetric := &models.MedicationMetric{
		Time: input.Time, Dosage: input.Dosage, Quantity: input.Quantity, Metric: input.Metric, Date: date, Name: input.Name}

	err = app.Models.MedicationMetric.InsertMedicationMetric(user.ID, medicationMetric)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.Background(func() {
		_ = app.Models.UserPoint.InsertPoint(user.ID, "Medication", 2)
	})
	env := envelope{
		"message": "Successfully Created User Medication Metrics!",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteMedicationMetric(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	err = app.Models.MedicationMetric.DeleteMedicationMetric(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Medication Metric successfully deleted"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
