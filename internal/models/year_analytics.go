package models

import (
	"context"
	"fmt"
	"time"
)

func (m AnalyticsModel) GetYearSymptomOccurrences(userID string, symptomID int, year int) (map[int]float64, error) {
	query := `
	SELECT
		EXTRACT(MONTH FROM date) AS month_of_year,
		AVG((morning_severity + afternoon_severity + night_severity) / 3) AS average_severity
	FROM
		user_symptoms_metric
	WHERE
		user_id = $1
		AND symptoms_id = $2
		AND EXTRACT(YEAR FROM date) = $3
	GROUP BY
		month_of_year
	ORDER BY
		month_of_year;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute the query with the provided parameters
	rows, err := m.DB.QueryContext(ctx, query, userID, symptomID, year)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	// Create a map to store symptom occurrences
	symptomOccurrences := make(map[int]float64)
	for rows.Next() {
		var sc SymptomCount
		// Scan the results into the SymptomCount struct
		err := rows.Scan(&sc.MonthOfYear, &sc.AvgSev)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		symptomOccurrences[sc.MonthOfYear] = Round(sc.AvgSev)
	}

	return symptomOccurrences, nil
}
func (m AnalyticsModel) GetYearBowelTypeOccurrences(userID string, year int) (map[int][]KeyValue, error) {
	query := `
	SELECT
		EXTRACT(MONTH FROM date) AS month_of_year,
		type,
		COUNT(*) AS occurrences
	FROM
		user_bowel_metric
	WHERE
		user_id = $1
		AND EXTRACT(YEAR FROM date) = $2
	GROUP BY
		month_of_year, type
	ORDER BY
		month_of_year, type;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID, year)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	bowelTypeOccurrences := make(map[int][]KeyValue)
	for rows.Next() {
		var monthOfYear, typeID, occurrences int
		err := rows.Scan(&monthOfYear, &typeID, &occurrences)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		bowelTypeOccurrences[monthOfYear] = append(bowelTypeOccurrences[monthOfYear], KeyValue{Key: typeID, Value: occurrences})
	}

	// Ensure all months from 1 to 12 have an entry in the map
	for i := 1; i <= 12; i++ {
		if _, exists := bowelTypeOccurrences[i]; !exists {
			bowelTypeOccurrences[i] = []KeyValue{}
		}
	}

	return bowelTypeOccurrences, nil
}
func (m AnalyticsModel) GetYearTagOccurrences(userID string, year int, tagToQuery string) (map[int][]KeyValue, error) {
	var query string

	if tagToQuery == "" {
		query = `
        WITH tag_occurrences AS (
            SELECT
                EXTRACT(MONTH FROM date) AS month_of_year,
                UNNEST(breakfast_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(YEAR FROM date) = $2
            UNION ALL
            SELECT
                EXTRACT(MONTH FROM date) AS month_of_year,
                UNNEST(lunch_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(YEAR FROM date) = $2
            UNION ALL
            SELECT
                EXTRACT(MONTH FROM date) AS month_of_year,
                UNNEST(dinner_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(YEAR FROM date) = $2
            UNION ALL
            SELECT
                EXTRACT(MONTH FROM date) AS month_of_year,
                UNNEST(snack_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(YEAR FROM date) = $2
        )
        SELECT
            month_of_year,
            tag,
            COUNT(*) AS occurrences
        FROM
            tag_occurrences
        GROUP BY
            month_of_year,
            tag
        ORDER BY
            month_of_year,
            tag;
        `
	} else {
		query = `
        WITH tag_occurrences AS (
            SELECT
                EXTRACT(MONTH FROM date) AS month_of_year,
                UNNEST(breakfast_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(YEAR FROM date) = $2
            UNION ALL
            SELECT
                EXTRACT(MONTH FROM date) AS month_of_year,
                UNNEST(lunch_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(YEAR FROM date) = $2
            UNION ALL
            SELECT
                EXTRACT(MONTH FROM date) AS month_of_year,
                UNNEST(dinner_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(YEAR FROM date) = $2
            UNION ALL
            SELECT
                EXTRACT(MONTH FROM date) AS month_of_year,
                UNNEST(snack_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND EXTRACT(YEAR FROM date) = $2
        )
        SELECT
            month_of_year,
            tag,
            COUNT(*) AS occurrences
        FROM
            tag_occurrences
        WHERE
            tag = $3
        GROUP BY
            month_of_year,
            tag
        ORDER BY
            month_of_year,
            tag;
        `
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{userID, year}
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
		var monthOfYear, occurrences int
		var tag string
		err := rows.Scan(&monthOfYear, &tag, &occurrences)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		tagOccurrences[monthOfYear] = append(tagOccurrences[monthOfYear], KeyValue{Key: tag, Value: occurrences})
	}

	// Ensure all months from 1 to 12 have an entry in the map
	for i := 1; i <= 12; i++ {
		if _, exists := tagOccurrences[i]; !exists {
			tagOccurrences[i] = []KeyValue{}
		}
	}

	return tagOccurrences, nil
}
