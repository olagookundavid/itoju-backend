package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/olagookundavid/itoju/cmd/api"
	"github.com/olagookundavid/itoju/internal/vcs"
	"github.com/robfig/cron/v3"
)

var (
	version = vcs.Version()
)

// func mainx() {
// 	app := intialize()
// 	intializeBackGroundTask(app)
// 	err := server.Serve(app)
// 	if err != nil {
// 		jsonlog.New(os.Stdout, jsonlog.LevelInfo).PrintFatal(err, nil)
// 	}
// }

func intializeBackGroundTask(app *api.Application) {
	app.Background(func() {
		cronJob(app)
	})
}

func loadDbUrl() string {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env variable missing")
	}
	return dbUrl
}

func expvarSetup(db *sql.DB) {
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
}

func displayVersion(flagStr string) {
	displayVersion := flag.Bool(flagStr, false, "Display version and exit")
	flag.Parse()
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}
}

func openDB(cfg api.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}
	// Setting no default connections
	// db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	// db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	// duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	// if err != nil {
	// 	return nil, err
	// }
	// db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func cronJob(app *api.Application) {
	// 	ticker := time.NewTicker(2 * time.Second)
	// 	defer ticker.Stop()
	//	for {
	//		select {
	//		case <-ticker.C:
	//			fmt.Println("Running task every second")
	//		}
	//	}
	c := cron.New()

	_, err := c.AddFunc("@daily", func() {
		app.Logger.PrintInfo("Deleting Tokens from tokens table", nil)
		err := app.Models.Tokens.DeleteAllExpiredTokens()
		if err != nil {
			app.Logger.PrintError(err, map[string]string{"error": "An error occured with deleting tokens from Tokens Table"})
			return
		}
	})

	if err != nil {
		app.Logger.PrintError(err, map[string]string{"error": "An error occured with the cron job"})
		return
	}

	c.Start()
}

func flagSetup(dbUrl string) *api.Config {

	var cfg api.Config

	//env and port
	flag.IntVar(&cfg.Port, "port", 8000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	//db
	flag.StringVar(&cfg.Db.Dsn, "db-dsn", dbUrl, "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 15, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 12, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", "1m", "PostgreSQL max connection idle time")
	//limiters
	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	return &cfg
}
