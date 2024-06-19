package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type BowelMetric struct {
	ID   int       `json:"id"`
	Time string    `json:"time"`
	Tags []string  `json:"tags"`
	Date time.Time `json:"date"`
	Type float64   `json:"type"`
	Pain float64   `json:"pain"`
}

type BowelMetricModel struct {
	DB *sql.DB
}

func (m BowelMetricModel) GetUserBowelMetrics(userId string, date time.Time) ([]*BowelMetric, error) {

	query := `
	SELECT ubm.id, ubm.time, ubm.type, ubm.pain, ubm.tags, ubm.date
    FROM user_bowel_metric ubm
    WHERE ubm.user_id = $1 AND ubm.date = $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bowelMetrics := []*BowelMetric{}
	for rows.Next() {
		var bowelMetric BowelMetric
		err := rows.Scan(&bowelMetric.ID, &bowelMetric.Time, &bowelMetric.Type, &bowelMetric.Pain, pq.Array(&bowelMetric.Tags), &bowelMetric.Date)
		if err != nil {
			return nil, err
		}

		bowelMetrics = append(bowelMetrics, &bowelMetric)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bowelMetrics, nil
}

func (m BowelMetricModel) GetUserBowelMetric(userId string, id int64) (*BowelMetric, error) {
	query := `
    SELECT ubm.id, ubm.time, ubm.type, ubm.pain, ubm.tags, ubm.date
    FROM user_bowel_metric ubm
    WHERE ubm.user_id = $1 AND ubm.id = $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, userId, id)

	var bowelMetric BowelMetric
	err := row.Scan(&bowelMetric.ID, &bowelMetric.Time, &bowelMetric.Type, &bowelMetric.Pain, pq.Array(&bowelMetric.Tags), &bowelMetric.Date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &bowelMetric, nil
}

func (m BowelMetricModel) InsertBowelMetric(userID string, bowelMetric *BowelMetric) error {

	query := `
	INSERT INTO user_bowel_metric (user_id, time, pain, type, date, tags)
	VALUES ($1, $2, $3, $4, $5, $6) `

	args := []any{userID, bowelMetric.Time, bowelMetric.Pain, bowelMetric.Type, bowelMetric.Date, pq.Array(bowelMetric.Tags)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m BowelMetricModel) UpdateBowelMetric(bowelMetric *BowelMetric) error {

	query := ` UPDATE user_bowel_metric SET time = $1, pain = $2, type = $3, tags = $4 WHERE id = $5; `

	args := []any{bowelMetric.Time, bowelMetric.Pain, bowelMetric.Type, pq.Array(bowelMetric.Tags), bowelMetric.ID}
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

func (m BowelMetricModel) DeleteBowelMetric(id int64, user_id string) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := ` DELETE FROM user_bowel_metric WHERE id = $1 AND user_id = $2 `
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

func (m BowelMetricModel) CheckUserEntry(userID string, date time.Time, sendbool chan<- bool) {

	query := `
	SELECT COUNT(*) AS entry_count
	FROM user_bowel_metric ubm
	WHERE ubm.user_id = $1 AND ubm.date = $2
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
