package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
type SymTopN struct {
	Id    int    `json:"id"`
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type SymsMetricModel struct {
	DB *sql.DB
}

func (m SymsMetricModel) CreateSymsMetric(userId string, symsMetric SymsMetric) error {
	query := `
	INSERT INTO user_symptoms_metric (user_id, symptoms_id, date)
	VALUES ($1, $2, $3) `

	args := []any{userId, symsMetric.Id, symsMetric.Date}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "unique_user_symptom_date"`:

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

func (m SymsMetricModel) GetUserTopNSyms(userId string, interval int) ([]*SymTopN, error) {

	query := fmt.Sprintf(
		`
	SELECT s.name, usm.symptoms_id, COUNT(*) AS count
	FROM user_symptoms_metric usm
	JOIN symptoms s ON usm.symptoms_id = s.id
	WHERE usm.user_id = $1
	AND usm.date >= CURRENT_DATE - INTERVAL '%d days'
	GROUP BY s.name, usm.symptoms_id
	ORDER BY count DESC
	LIMIT 4; 
	`, interval)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outPuts := []*SymTopN{}
	for rows.Next() {
		var outPut SymTopN
		err := rows.Scan(&outPut.Name, &outPut.Id, &outPut.Count)
		if err != nil {
			return nil, err
		}

		outPuts = append(outPuts, &outPut)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return outPuts, nil
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

func (m SymsMetricModel) DaysTrackedInARow(userID string) (*int, error) {

	query := `
        SELECT MAX(consecutive_days) AS max_consecutive_days
        FROM (
            SELECT COUNT(*) AS consecutive_days
            FROM (
                SELECT date,
                       ROW_NUMBER() OVER (ORDER BY date) - 
                       ROW_NUMBER() OVER (PARTITION BY tracked ORDER BY date) AS grp
                FROM (
                    SELECT date, 
                           CASE WHEN COUNT(user_id) > 0 THEN 1 ELSE 0 END AS tracked
                    FROM user_symptoms_metric
                    WHERE user_id = $1
                    GROUP BY date
                ) AS t
            ) AS s
            GROUP BY grp
        ) AS max_consecutive_days_query
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var maxConsecutiveDays int
	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&maxConsecutiveDays)
	if err != nil {
		return nil, err
	}
	return &maxConsecutiveDays, err
}

func (m SymsMetricModel) DaysTrackedFree(userID string) (*int, error) {

	query := `
	SELECT COUNT(*) AS max_consecutive_symptom_free_days
	FROM (
		SELECT date,
			   CASE WHEN LAG(tracked, 1, 0) OVER (ORDER BY date) = 0 THEN 1 ELSE 0 END AS is_consecutive
		FROM (
			SELECT g.date, 
				   CASE WHEN COUNT(usm.user_id) = 0 THEN 1 ELSE 0 END AS tracked
			FROM generate_series(CURRENT_DATE, CURRENT_DATE - INTERVAL '29 days', '-1 day') AS g(date)
			LEFT JOIN user_symptoms_metric AS usm ON g.date = usm.date AND usm.user_id = $1
			GROUP BY g.date
		) AS t
	) AS s
	WHERE is_consecutive = 1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var maxConsecutiveDays int
	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&maxConsecutiveDays)
	if err != nil {
		return nil, err
	}
	return &maxConsecutiveDays, err
}

func (m SymsMetricModel) CheckUserEntry(userID string, date time.Time, sendbool chan<- bool) {

	query := `
	SELECT COUNT(*) AS entry_count
	FROM user_symptoms_metric usm
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
