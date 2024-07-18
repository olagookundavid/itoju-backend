package api

import (
	"net/http"
)

func (app *Application) GetMonthBowelAnalytics(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	month, err := app.readIntParam(r, "month")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	analytics, err := app.Models.AnalyticsMetric.GetMonthBowelTypeOccurrences(user.ID, int(month))
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

func (app *Application) GetSymsMonthAnalytics(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}
	month, err := app.readIntParam(r, "month")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	analytics, err := app.Models.AnalyticsMetric.GetMonthSymptomOccurrences(user.ID, int(id), int(month))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":          "Retrieved All Analytics for user",
		"analyticsMetrics": EnsureAllWeeksPresent(analytics)}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) GetTagsMonthAnalytics(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	month, err := app.readIntParam(r, "month")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	tagToQuery, err := app.readStringParam(r, "tag")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if tagToQuery == "General" {
		tagToQuery = ""
	}

	analytics, err := app.Models.AnalyticsMetric.GetMonthTagOccurrences(user.ID, int(month), tagToQuery)
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

func EnsureAllWeeksPresent(metrics map[int]float64) map[int]float64 {
	for i := 0; i <= 5; i++ {
		if _, exists := metrics[i]; !exists {
			metrics[i] = 0
		}
	}
	return metrics
}
