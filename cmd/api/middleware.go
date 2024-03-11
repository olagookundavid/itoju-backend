package api

import (
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
	"github.com/olagookundavid/itoju/internal/validator"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *Application) RecoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			// Use the builtin recover function to check if there has been a panic or // not.
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.Config.Limiter.Enabled {
			println("rater")
			ip := realip.FromRequest(r)
			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(app.Config.Limiter.Rps),
						app.Config.Limiter.Burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				println("rater2")
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			println("rater3")
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}
	})
}

func (app *Application) RequireActivatedAndAuthedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		status := app.contextGetStatus(r)

		if !status {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetUser(r, models.AnonymousUser)
			r = app.contextSetTokenStatus(r, true)
			next.ServeHTTP(w, r)
			println("anon1")
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			r = app.contextSetUser(r, models.AnonymousUser)
			r = app.contextSetTokenStatus(r, false)
			next.ServeHTTP(w, r)
			println("anon2")
			return
		}
		token := headerParts[1]
		v := validator.New()
		if models.ValidateTokenPlaintext(v, token); !v.Valid() {
			r = app.contextSetUser(r, models.AnonymousUser)
			r = app.contextSetTokenStatus(r, false)
			next.ServeHTTP(w, r)
			println("anon2")
			return
		}

		user, err := app.Models.Users.GetForToken(models.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				r = app.contextSetUser(r, models.AnonymousUser)
				r = app.contextSetTokenStatus(r, false)
				next.ServeHTTP(w, r)
				println("anon3")
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetTokenStatus(r, true)
		r = app.contextSetUser(r, user)
		println("not anon")
		next.ServeHTTP(w, r)
	})
}

func (app *Application) Metrics(next http.Handler) http.Handler {
	var (
		totalRequestsReceived           = expvar.NewInt("total_requests_received")
		totalResponsesSent              = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
		totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")
	)
	// The following code will be run for every request...
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the time that we started to process the request.
		start := time.Now()
		totalRequestsReceived.Add(1) // Call the next handler in the chain.
		mw := &metricsResponseWriter{wrapped: w}
		next.ServeHTTP(mw, r)
		totalResponsesSent.Add(1)
		totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}

type metricsResponseWriter struct {
	wrapped       http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func (mw *metricsResponseWriter) Header() http.Header {
	return mw.wrapped.Header()
}
func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)
	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}
func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	if !mw.headerWritten {
		mw.statusCode = http.StatusOK
		mw.headerWritten = true
	}
	return mw.wrapped.Write(b)
}
func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped
}
