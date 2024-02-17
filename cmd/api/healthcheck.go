package api

import (
	"net/http"

	"github.com/olagookundavid/itoju/internal/vcs"
)

func (app *Application) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{"status": "available",
		"system_info": map[string]string{"environment": app.Config.Env,
			"version": vcs.Version()}}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
