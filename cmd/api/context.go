package api

import (
	"context"
	"net/http"

	"github.com/olagookundavid/itoju/internal/models"
)

type contextKey string

const userContextKey = contextKey("user")
const statusContextKey = contextKey("status")

func (app *Application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *Application) contextSetTokenStatus(r *http.Request, status bool) *http.Request {
	ctx := context.WithValue(r.Context(), statusContextKey, status)
	return r.WithContext(ctx)
}

func (app *Application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}

func (app *Application) contextGetStatus(r *http.Request) bool {
	status, ok := r.Context().Value(statusContextKey).(bool)
	if !ok {
		panic("missing status value")
	}
	return status
}
