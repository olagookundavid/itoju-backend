package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type MedicationMetric struct {
	ID       int       `json:"id"`
	Time     string    `json:"time"`
	Name     string    `json:"name"`
	Metric   string    `json:"metric"`
	Date     time.Time `json:"date"`
	Dosage   float64   `json:"dosage"`
	Quantity float64   `json:"quantity"`
}

type MedicationMetricModel struct {
	DB *sql.DB
}

func (m MedicationMetricModel) GetUserMedicationMetrics(userId string, date time.Time) ([]*MedicationMetric, error) {

	query := `
	SELECT umm.id, umm.time, umm.dosage, umm.quantity, umm.name, umm.date, umm.metric
    FROM user_medication_metric umm
    WHERE umm.user_id = $1 AND umm.date = $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	medicationMetrics := []*MedicationMetric{}
	for rows.Next() {
		var medicationMetric MedicationMetric
		err := rows.Scan(&medicationMetric.ID, &medicationMetric.Time, &medicationMetric.Dosage, &medicationMetric.Quantity, &medicationMetric.Name, &medicationMetric.Date, &medicationMetric.Metric)
		if err != nil {
			return nil, err
		}

		medicationMetrics = append(medicationMetrics, &medicationMetric)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return medicationMetrics, nil
}

func (m MedicationMetricModel) GetUserMedicationMetric(userId string, id int64) (*MedicationMetric, error) {
	query := `
    SELECT umm.id, umm.time, umm.dosage, umm.quantity, umm.name, umm.date, umm.metric
    FROM user_medication_metric umm
    WHERE umm.user_id = $1 AND umm.id = $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, userId, id)

	var medicationMetric MedicationMetric
	err := row.Scan(&medicationMetric.ID, &medicationMetric.Time, &medicationMetric.Dosage, &medicationMetric.Quantity, &medicationMetric.Name, &medicationMetric.Date, &medicationMetric.Metric)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &medicationMetric, nil
}

func (m MedicationMetricModel) InsertMedicationMetric(userID string, medicationMetric *MedicationMetric) error {

	query := `
	INSERT INTO user_medication_metric (user_id, time, dosage, quantity, date, name, metric)
	VALUES ($1, $2, $3, $4, $5, $6, $7) `

	args := []any{userID, medicationMetric.Time, medicationMetric.Dosage, medicationMetric.Quantity, medicationMetric.Date, medicationMetric.Name, medicationMetric.Metric}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m MedicationMetricModel) UpdateMedicationMetric(medicationMetric *MedicationMetric) error {

	query := ` UPDATE user_medication_metric SET time = $1, dosage = $2, quantity = $3, metric = $4, name = $5 WHERE id = $6; `

	args := []any{medicationMetric.Time, medicationMetric.Dosage, medicationMetric.Quantity, medicationMetric.Metric, medicationMetric.Name, medicationMetric.ID}
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

func (m MedicationMetricModel) DeleteMedicationMetric(id int64, user_id string) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := ` DELETE FROM user_medication_metric WHERE id = $1 AND user_id = $2 `
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

func (m MedicationMetricModel) CheckUserEntry(userID string, date time.Time, sendbool chan<- bool) {

	query := `
	SELECT COUNT(*) AS entry_count
	FROM user_medication_metric umm
	WHERE umm.user_id = $1 AND umm.date = $2
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
