package api

import (
	"net/http"
)

func (app *Application) GetBowelDaysAnalytics(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	days, err := app.readIntParam(r, "days")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	analytics, err := app.Models.AnalyticsMetric.GetBowelTypeOccurrences(user.ID, int(days))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":          "Retrieved All Analytics for user",
		"analyticsMetrics": analytics}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) GetSymsDaysAnalytics(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}
	days, err := app.readIntParam(r, "days")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	analytics, err := app.Models.AnalyticsMetric.GetSymptomOccurrences(user.ID, int(id), int(days))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":          "Retrieved All Analytics for user",
		"analyticsMetrics": EnsureAllDaysPresent(analytics)}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) GetTagsDaysAnalytics(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	days, err := app.readIntParam(r, "days")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	analytics, err := app.Models.AnalyticsMetric.GetTagOccurrences(user.ID, int(days), "")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":          "Retrieved All Analytics for user",
		"analyticsMetrics": analytics}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func EnsureAllDaysPresent(metrics map[int]int) map[int]int {
	for i := 0; i <= 6; i++ {
		if _, exists := metrics[i]; !exists {
			metrics[i] = 0
		}
	}
	return metrics
}
