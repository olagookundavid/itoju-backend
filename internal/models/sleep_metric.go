package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type SleepMetric struct {
	ID         int       `json:"id"`
	IsNight    bool      `json:"is_night"`
	TimeSlept  string    `json:"time_slept"`
	TimeWokeUp string    `json:"time_woke_up"`
	Tags       []string  `json:"tags"`
	Date       time.Time `json:"date"`
	Severity   float64   `json:"severity"`
}

type SleepMetricModel struct {
	DB *sql.DB
}

func (m SleepMetricModel) GetUserSleepMetric(userId string, date time.Time, isNight bool) (*SleepMetric, error) {
	query := `
    SELECT usm.id, usm.is_night, usm.time_slept, usm.time_woke_up, usm.tags, usm.date, usm.severity
    FROM user_sleep_metric usm
    WHERE usm.user_id = $1 AND usm.date = $2 AND usm.is_night = $3
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, userId, date, isNight)

	var sleepMetric SleepMetric
	err := row.Scan(&sleepMetric.ID, &sleepMetric.IsNight, &sleepMetric.TimeSlept, &sleepMetric.TimeWokeUp, pq.Array(&sleepMetric.Tags), &sleepMetric.Date, &sleepMetric.Severity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &sleepMetric, nil
}

func (m SleepMetricModel) InsertSleepMetric(userID string, sleepMetric *SleepMetric) error {
	uery := `
	CREATE TABLE user_sleep_metric (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
		is_night BOOLEAN NOT NULL,
		time_slept VARCHAR(20) NOT NULL,
		time_woke_up VARCHAR(20) NOT NULL, 
		tags TEXT[] NOT NULL,
		date DATE NOT NULL DEFAULT CURRENT_DATE,
		severity NUMERIC(3,2) DEFAULT 0 CHECK (severity >= 0 AND severity <= 1),
		CONSTRAINT unique_user_sleep_date UNIQUE (user_id, is_night, date)
	); `

	query := `
	INSERT INTO user_sleep_metric (user_id, is_night, time_slept, time_woke_up, date, severity, tags)
	VALUES ($1, $2, $3, $4, $5, $6, $7) `

	args := []any{userID, sleepMetric.IsNight, sleepMetric.TimeSlept, sleepMetric.TimeWokeUp, sleepMetric.Date, sleepMetric.Severity, pq.Array(sleepMetric.Tags)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, _ = m.DB.ExecContext(ctx, uery)
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m SleepMetricModel) UpdateSleepMetric(sleepMetric *SleepMetric) error {

	query := ` UPDATE user_sleep_metric SET time_slept = $1, time_woke_up = $2, tags = $3, severity = $4 WHERE id = $5 AND is_night = $6; `

	args := []any{sleepMetric.TimeSlept, sleepMetric.TimeWokeUp, pq.Array(sleepMetric.Tags), sleepMetric.Severity, sleepMetric.ID, sleepMetric.IsNight}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

/*
http://localhost:8000/v1/user/sleep_metrics/2024-03-3/true


*/