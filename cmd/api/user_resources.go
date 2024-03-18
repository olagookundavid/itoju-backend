package api

import (
	"errors"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
	"github.com/olagookundavid/itoju/internal/validator"
)

func (app *Application) GetResources(w http.ResponseWriter, r *http.Request) {

	resources, err := app.Models.Resources.GetResources()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"message":   "Retrieved All Resources",
		"resources": resources}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) InsertResources(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string   `json:"name"`
		ImageUrl string   `json:"image_url"`
		Link     string   `json:"link"`
		Tags     []string `json:"tags"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	resource := &models.Resources{
		Name: input.Name, Tags: input.Tags, ImageUrl: input.ImageUrl, Link: input.Link,
	}

	err = app.Models.Resources.InsertResources(*resource)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordAlreadyExist):
			app.recordAlreadyExistsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	env := envelope{
		"message": "Successfully added Resource",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) UpdateResources(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	resource, err := app.Models.Resources.Get(id)
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
		Name     *string  `json:"name"`
		ImageUrl *string  `json:"image_url"`
		Link     *string  `json:"link"`
		Tags     []string `json:"tags"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		resource.Name = *input.Name
	}
	if input.ImageUrl != nil {
		resource.ImageUrl = *input.ImageUrl
	}
	if input.Link != nil {
		resource.Link = *input.Link
	}
	if input.Tags != nil {
		resource.Tags = input.Tags
	}

	v := validator.New()
	if models.ValidateResource(v, resource); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.Models.Resources.UpdateResources(resource)
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
		"message":  "Successfully updated Resource",
		"resource": resource,
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteResources(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}
	err = app.Models.Resources.DeleteResources(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Resource successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
