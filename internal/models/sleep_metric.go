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

func (m SleepMetricModel) GetUserSleepMetrics(userId string, date time.Time) ([]*SleepMetric, error) {

	query := `
	SELECT usm.id, usm.is_night, usm.time_slept, usm.time_woke_up, usm.tags, usm.date, usm.severity
    FROM user_sleep_metric usm
    WHERE usm.user_id = $1 AND usm.date = $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sleepsMetrics := []*SleepMetric{}
	for rows.Next() {
		var sleepMetric SleepMetric
		err := rows.Scan(&sleepMetric.ID, &sleepMetric.IsNight, &sleepMetric.TimeSlept, &sleepMetric.TimeWokeUp, pq.Array(&sleepMetric.Tags), &sleepMetric.Date, &sleepMetric.Severity)
		if err != nil {
			return nil, err
		}

		sleepsMetrics = append(sleepsMetrics, &sleepMetric)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sleepsMetrics, nil
}

func (m SleepMetricModel) GetUserSleepMetric(userId string, id int64) (*SleepMetric, error) {
	query := `
    SELECT usm.id, usm.is_night, usm.time_slept, usm.time_woke_up, usm.tags, usm.date, usm.severity
    FROM user_sleep_metric usm
    WHERE usm.user_id = $1 AND usm.id = $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, userId, id)

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

	query := `
	INSERT INTO user_sleep_metric (user_id, is_night, time_slept, time_woke_up, date, severity, tags)
	VALUES ($1, $2, $3, $4, $5, $6, $7) `

	args := []any{userID, sleepMetric.IsNight, sleepMetric.TimeSlept, sleepMetric.TimeWokeUp, sleepMetric.Date, sleepMetric.Severity, pq.Array(sleepMetric.Tags)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
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

func (m SleepMetricModel) DeleteSleepMetric(id int64, user_id string) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := ` DELETE FROM user_sleep_metric WHERE id = $1 AND user_id = $2 `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id, user_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m SleepMetricModel) CheckUserEntry(userID string, date time.Time, sendbool chan<- bool) {

	query := `
	SELECT COUNT(*) AS entry_count
	FROM user_sleep_metric usm
	WHERE usm.user_id = $1 AND usm.date = $2
`
	var entryCount int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, userID, date).Scan(&entryCount)
	if err != nil {
		sendbool <- false
		return
	}
	sendbool <- (entryCount > 0)
}
