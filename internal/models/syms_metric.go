package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type SymsMetric struct {
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	Date              time.Time `json:"date"`
	MorningSeverity   float32   `json:"morning_severity"`
	AfternoonSeverity float32   `json:"afternoon_severity"`
	NightSeverity     float32   `json:"night_severity"`
}

type SymsMetricModel struct {
	DB *sql.DB
}

func (m SymsMetricModel) CreateSymsMetric(userId string, symsMetric SymsMetric) error {
	query := `
	INSERT INTO user_symptoms_metric (user_id, symptoms_id)
	VALUES ($1, $2) `

	args := []any{userId, symsMetric.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "unique_user_symptom_date_us"`:

			return ErrRecordAlreadyExist
		default:
			return err
		}
	}
	return err
}

func (m SymsMetricModel) GetUserSymptomsMetric(userId string, date time.Time) ([]*SymsMetric, error) {
	query := `
	SELECT usm.id, s.name, usm.date, usm.morning_severity, usm.afternoon_severity, usm.night_severity
	FROM user_symptoms_metric usm
	JOIN symptoms s ON usm.symptoms_id = s.id
	WHERE usm.user_id = $1 AND usm.date = $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	symsMetrics := []*SymsMetric{}
	for rows.Next() {
		var symsMetric SymsMetric
		err := rows.Scan(&symsMetric.Id, &symsMetric.Name, &symsMetric.Date, &symsMetric.MorningSeverity, &symsMetric.AfternoonSeverity, &symsMetric.NightSeverity)
		if err != nil {
			return nil, err
		}

		symsMetrics = append(symsMetrics, &symsMetric)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return symsMetrics, nil
}

func (m SymsMetricModel) Get(id int64) (*SymsMetric, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := ` SELECT usm.id, usm.date, usm.morning_severity, usm.afternoon_severity, usm.night_severity
	FROM user_symptoms_metric usm WHERE id = $1; `
	var symsMetric SymsMetric
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&symsMetric.Id, &symsMetric.Date, &symsMetric.MorningSeverity, &symsMetric.AfternoonSeverity, &symsMetric.NightSeverity)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &symsMetric, nil
}

func (m SymsMetricModel) UpdateSymsMetric(symsMetric *SymsMetric, id int) error {

	query := ` UPDATE user_symptoms_metric SET morning_severity = $1, afternoon_severity= $2, night_severity = $3 WHERE id = $4; `

	args := []any{symsMetric.MorningSeverity, symsMetric.AfternoonSeverity, symsMetric.NightSeverity, symsMetric.Id}
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

func (m SymsMetricModel) DeleteSymsMetric(id int64, user_id string) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := ` DELETE FROM user_symptoms_metric WHERE id = $1 AND user_id = $2 `
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
