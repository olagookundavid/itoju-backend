package models

import (
	"context"
	"fmt"
	"time"
)

func (m AnalyticsModel) GetMonthSymptomOccurrences(userID string, symptomID int, month int) (map[int]float64, error) {
	query := `
	SELECT
		(EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
		AVG((morning_severity + afternoon_severity + night_severity) / 3) AS average_severity
	FROM
		user_symptoms_metric
	WHERE
		user_id = $1
		AND symptoms_id = $2
		AND EXTRACT(MONTH FROM date) = $3
	GROUP BY
		week_of_month
	ORDER BY
		week_of_month;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute the query with the provided parameters
	rows, err := m.DB.QueryContext(ctx, query, userID, symptomID, month)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	// Create a map to store symptom occurrences
	symptomOccurrences := make(map[int]float64)
	for rows.Next() {
		var sc SymptomCount
		// Scan the results into the SymptomCount struct
		err := rows.Scan(&sc.WeekOfMonth, &sc.AvgSev)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		symptomOccurrences[sc.WeekOfMonth] = Round(sc.AvgSev)
	}

	return symptomOccurrences, nil
}

func (m AnalyticsModel) GetMonthBowelTypeOccurrences(userID string, month int) (map[int][]KeyValue, error) {
	query := `
		SELECT
			(EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
			type,
			COUNT(*) AS occurrences
		FROM
			user_bowel_metric
		WHERE
			user_id = $1
			AND EXTRACT(MONTH FROM date) = $2
		GROUP BY
			week_of_month, type
		ORDER BY
			week_of_month, type;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID, month)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	bowelTypeOccurrences := make(map[int][]KeyValue)
	for rows.Next() {
		var weekOfMonth, typeID, occurrences int
		err := rows.Scan(&weekOfMonth, &typeID, &occurrences)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		bowelTypeOccurrences[weekOfMonth] = append(bowelTypeOccurrences[weekOfMonth], KeyValue{Key: typeID, Value: occurrences})
	}
	for i := 0; i <= 6; i++ {
		if _, exists := bowelTypeOccurrences[i]; !exists {
			bowelTypeOccurrences[i] = []KeyValue{}
			bowelTypeOccurrences[i] = []KeyValue{}
		}
	}

	return bowelTypeOccurrences, nil
}

func (m AnalyticsModel) GetMonthTagOccurrences(userID string, month int, tagToQuery string) (map[int][]KeyValue, error) {
	var query string

	if tagToQuery == "" {
		query = `
        WITH tag_occurrences AS (
            SELECT
                (EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
                UNNEST(breakfast_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(MONTH FROM date) = $2
            UNION ALL
            SELECT
                (EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
                UNNEST(lunch_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(MONTH FROM date) = $2
            UNION ALL
            SELECT
                (EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
                UNNEST(dinner_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(MONTH FROM date) = $2
            UNION ALL
            SELECT
                (EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
                UNNEST(snack_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(MONTH FROM date) = $2
        )
        SELECT
            week_of_month,
            tag,
            COUNT(*) AS occurrences
        FROM
            tag_occurrences
        GROUP BY
            week_of_month,
            tag
        ORDER BY
            week_of_month,
            tag;
		`
	} else {
		query = `
        WITH tag_occurrences AS (
            SELECT
                (EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
                UNNEST(breakfast_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(MONTH FROM date) = $2
            UNION ALL
            SELECT
                (EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
                UNNEST(lunch_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(MONTH FROM date) = $2
            UNION ALL
            SELECT
                (EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
                UNNEST(dinner_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(MONTH FROM date) = $2
            UNION ALL
            SELECT
                (EXTRACT(WEEK FROM date) - EXTRACT(WEEK FROM date_trunc('month', date)) + 1) AS week_of_month,
                UNNEST(snack_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(MONTH FROM date) = $2
        )
        SELECT
            week_of_month,
            tag,
            COUNT(*) AS occurrences
        FROM
            tag_occurrences
        WHERE
            tag = $3
        GROUP BY
            week_of_month,
            tag
        ORDER BY
            week_of_month,
            tag;
        `
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{userID, month}
	if tagToQuery != "" {
		args = append(args, tagToQuery)
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	tagOccurrences := make(map[int][]KeyValue)
	for rows.Next() {
		var weekOfMonth, occurrences int
		var tag string
		err := rows.Scan(&weekOfMonth, &tag, &occurrences)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		tagOccurrences[weekOfMonth] = append(tagOccurrences[weekOfMonth], KeyValue{Key: tag, Value: occurrences})
	}
	for i := 0; i <= 5; i++ {
		if _, exists := tagOccurrences[i]; !exists {
			tagOccurrences[i] = []KeyValue{}
		}
	}

	return tagOccurrences, nil
}
