package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/olagookundavid/itoju/internal/jsonlog"
)

type Metrics struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type MetricsModel struct {
	DB *sql.DB
}

// unique issue on user_metrics table, also how to do this well
func (m MetricsModel) SetUserMetrics(selectedMetrics []int, userID string) error {
	wg := sync.WaitGroup{}
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	query := ` INSERT INTO user_trackedmetric (user_id, metric_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, metricID := range selectedMetrics {
		wg.Add(1)
		go func(metricID int) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					logger.PrintError(fmt.Errorf("%s", err), nil)
				}
			}()
			_, err := m.DB.ExecContext(ctx, query, userID, metricID)
			if err != nil {
				return
			}
		}(metricID)
	}
	wg.Wait()
	return nil

}

func (m MetricsModel) GetMetrics() ([]*Metrics, error) {
	query := ` SELECT * FROM trackedmetrics `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return []*Metrics{}, err
	}
	defer rows.Close()
	metrics := []*Metrics{}
	for rows.Next() {
		var metric Metrics
		err := rows.Scan(&metric.Id, &metric.Name)
		if err != nil {
			return []*Metrics{}, err
		}

		metrics = append(metrics, &metric)
	}
	if err = rows.Err(); err != nil {
		return []*Metrics{}, err
	}
	return metrics, nil
}

func (m MetricsModel) GetUserMetrics(userID string) ([]*Metrics, error) {

	query := ` SELECT trackedmetrics.id , trackedmetrics.name FROM trackedmetrics
	JOIN user_trackedmetric ON trackedmetrics.id = user_trackedmetric.metric_id
	WHERE user_trackedmetric.user_id = $1; `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	metrics := []*Metrics{}
	for rows.Next() {
		var metric Metrics
		err := rows.Scan(&metric.Id, &metric.Name)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return metrics, nil
}

func (m MetricsModel) DeleteUserMetrics(userId string, selectedMetrics []int) error {
	wg := sync.WaitGroup{}
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	query := ` DELETE FROM user_trackedmetric
	WHERE user_id = $1
	AND metric_id = $2; `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, metricID := range selectedMetrics {

		wg.Add(1)
		go func(metricID int) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					logger.PrintError(fmt.Errorf("%s", err), nil)
				}
			}()
			_, err := m.DB.ExecContext(ctx, query, userId, metricID)
			if err != nil {
				return
			}
		}(metricID)
	}

	wg.Wait()
	return nil
}
