package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type ExerciseMetric struct {
	ID        int       `json:"id"`
	UserID    string    `json:"-"`
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
	Started   string    `json:"start"`
	Ended     string    `json:"ended"`
	Tags      []string  `json:"tags"`
	NoOfTimes int       `json:"no_of_times"`
}

type ExerciseMetricModel struct {
	DB *sql.DB
}

func (m ExerciseMetricModel) InsertExerciseMetric(exerciseMetric *ExerciseMetric) error {
	query := `
        INSERT INTO user_exercise_metric (user_id, date, name)
        VALUES ($1, $2, $3)
    `
	args := []any{
		exerciseMetric.UserID,
		exerciseMetric.Date, exerciseMetric.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m ExerciseMetricModel) GetUserExerciseMetric(userId string, date time.Time) ([]*ExerciseMetric, error) {
	query := `
    SELECT uem.id, uem.name, uem.started, uem.ended, uem.tags, uem.date, uem.no_of_times
    FROM user_exercise_metric uem
    WHERE uem.user_id = $1 AND uem.date = $2
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exerciseMetrics := []*ExerciseMetric{}
	for rows.Next() {
		var exerciseMetric ExerciseMetric
		err := rows.Scan(&exerciseMetric.ID, &exerciseMetric.Name, &exerciseMetric.Started, &exerciseMetric.Ended, pq.Array(&exerciseMetric.Tags), &exerciseMetric.Date, &exerciseMetric.NoOfTimes)
		if err != nil {
			return nil, err
		}

		exerciseMetrics = append(exerciseMetrics, &exerciseMetric)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return exerciseMetrics, nil
}

func (m ExerciseMetricModel) UpdateExerciseMetric(exerciseMetric *ExerciseMetric, id int) error {

	query := ` UPDATE user_exercise_metric SET started = $1, ended = $2, tags = $3, no_of_times = $4 WHERE id = $5; `

	args := []any{&exerciseMetric.Started, &exerciseMetric.Ended, pq.Array(&exerciseMetric.Tags), &exerciseMetric.NoOfTimes, id}
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

func (m ExerciseMetricModel) Get(id int64) (*ExerciseMetric, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
    SELECT uem.id, uem.name, uem.started, uem.ended, uem.tags, uem.date, uem.no_of_times
    FROM user_exercise_metric uem
    WHERE id = $1
    `

	var exerciseMetric ExerciseMetric
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exerciseMetric.ID, &exerciseMetric.Name, &exerciseMetric.Started, &exerciseMetric.Ended, pq.Array(&exerciseMetric.Tags), &exerciseMetric.Date, &exerciseMetric.NoOfTimes)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &exerciseMetric, nil
}

func (m ExerciseMetricModel) CheckUserEntry(userID string, date time.Time, sendbool chan<- bool) {

	query := `
	SELECT COUNT(*) AS entry_count
	FROM user_exercise_metric uem
	WHERE uem.user_id = $1 AND uem.date = $2
`
	var entryCount int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, userID, date).Scan(&entryCount)
	if err != nil {
		sendbool <- false
	}
	sendbool <- (entryCount > 0)
}
