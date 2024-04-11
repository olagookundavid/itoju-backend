package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/olagookundavid/itoju/cmd/api"
	"github.com/olagookundavid/itoju/internal/jsonlog"
	"github.com/olagookundavid/itoju/internal/routes"
)

func Serve(app *api.Application) error {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      routes.Routes(app),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownError := make(chan error)

	app.Background(
		func() {
			// Intercept the signals, as before.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			s := <-quit
			app.Logger.PrintInfo("shutting down server", map[string]string{"signal": s.String()})
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err := srv.Shutdown(ctx)
			if err != nil {
				shutdownError <- err
			}
			app.Logger.PrintInfo("completing background tasks", map[string]string{"addr": srv.Addr})

			shutdownError <- nil
		})

	logger.PrintInfo("starting server", map[string]string{"addr": srv.Addr, "env": app.Config.Env})
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	app.Logger.PrintInfo("stopped server", map[string]string{"addr": srv.Addr})
	return nil

}
