package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Analytics struct {
}
type AnalyticsModel struct {
	DB *sql.DB
}

// getSymptomOccurrences retrieves the count of symptom occurrences for the specified period
func (m AnalyticsModel) GetSymptomOccurrences(userID string, symptomID int, days int) (map[int]int, error) {
	query := fmt.Sprintf(`
	SELECT
		EXTRACT(DOW FROM date) AS day_of_week,
		COUNT(*) AS occurrences
	FROM
		user_symptoms_metric
	WHERE
		user_id = $1
		AND symptoms_id = $2
		AND date >= CURRENT_DATE - INTERVAL '%d days'
	GROUP BY
		day_of_week
	ORDER BY
		day_of_week;
`, days)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID, symptomID)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	symptomOccurrences := make(map[int]int)
	for rows.Next() {
		var sc SymptomCount
		err := rows.Scan(&sc.DayOfWeek, &sc.Occurrences)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		symptomOccurrences[sc.DayOfWeek] = sc.Occurrences
	}

	return symptomOccurrences, nil
}

type SymptomCount struct {
	DayOfWeek   int `json:"day_of_week"`
	Occurrences int `json:"occurrences"`
}

func (m AnalyticsModel) GetBowelTypeOccurrences(userID string, days int) (map[string][]KeyValue, error) {
	query := fmt.Sprintf(`
		SELECT
			EXTRACT(DOW FROM date) AS day_of_week,
			type,
			COUNT(*) AS occurrences
		FROM
			user_bowel_metric
		WHERE
			user_id = $1
			AND date >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY
			day_of_week, type
		ORDER BY
			day_of_week, type;
	`, days)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	bowelTypeOccurrences := make(map[string][]KeyValue)
	for rows.Next() {
		var dayOfWeek, typeID, occurrences int
		err := rows.Scan(&dayOfWeek, &typeID, &occurrences)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		bowelTypeOccurrences[strconv.Itoa(dayOfWeek)] = append(bowelTypeOccurrences[strconv.Itoa(dayOfWeek)], KeyValue{Key: typeID, Value: occurrences})
	}
	for i := 0; i <= 6; i++ {
		if _, exists := bowelTypeOccurrences[strconv.Itoa(i)]; !exists {
			bowelTypeOccurrences[strconv.Itoa(i)] = []KeyValue{}
		}
	}

	return bowelTypeOccurrences, nil
}

type KeyValue struct {
	Key   interface{} `json:"key"`
	Value int         `json:"value"`
}

func (m AnalyticsModel) GetTagOccurrences(userID string, days int, tagToQuery string) (map[string][]KeyValue, error) {
	var query string

	if tagToQuery == "" {
		query = fmt.Sprintf(`
        WITH tag_occurrences AS (
            SELECT
                EXTRACT(DOW FROM date) AS day_of_week,
                UNNEST(breakfast_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND date >= CURRENT_DATE - INTERVAL '%d days'
            UNION ALL
            SELECT
                EXTRACT(DOW FROM date) AS day_of_week,
                UNNEST(lunch_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND date >= CURRENT_DATE - INTERVAL '%d days'
            UNION ALL
            SELECT
                EXTRACT(DOW FROM date) AS day_of_week,
                UNNEST(dinner_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND date >= CURRENT_DATE - INTERVAL '%d days'
            UNION ALL
            SELECT
                EXTRACT(DOW FROM date) AS day_of_week,
                UNNEST(snack_tags) AS tag
            FROM
                user_food_metric
            WHERE
                user_id = $1
                AND date >= CURRENT_DATE - INTERVAL '%d days'
        )
        SELECT
            day_of_week,
            tag,
            COUNT(*) AS occurrences
        FROM
            tag_occurrences
        GROUP BY
            day_of_week,
            tag
        ORDER BY
            day_of_week,
            tag;
    `, days, days, days, days)
	} else {

		query = fmt.Sprintf(`
	WITH tag_occurrences AS (
		SELECT
			EXTRACT(DOW FROM date) AS day_of_week,
			UNNEST(breakfast_tags) AS tag
		FROM
			user_food_metric
		WHERE
			user_id = $1
			AND date >= CURRENT_DATE - INTERVAL '%d days'
		UNION ALL
		SELECT
			EXTRACT(DOW FROM date) AS day_of_week,
			UNNEST(lunch_tags) AS tag
		FROM
			user_food_metric
		WHERE
			user_id = $1
			AND date >= CURRENT_DATE - INTERVAL '%d days'
		UNION ALL
		SELECT
			EXTRACT(DOW FROM date) AS day_of_week,
			UNNEST(dinner_tags) AS tag
		FROM
			user_food_metric
		WHERE
			user_id = $1
			AND date >= CURRENT_DATE - INTERVAL '%d days'
		UNION ALL
		SELECT
			EXTRACT(DOW FROM date) AS day_of_week,
			UNNEST(snack_tags) AS tag
		FROM
			user_food_metric
		WHERE
			user_id = $1
			AND date >= CURRENT_DATE - INTERVAL '%d days'
	)
	SELECT
		day_of_week,
		tag,
		COUNT(*) AS occurrences
	FROM
		tag_occurrences
	WHERE
		tag = $2
	GROUP BY
		day_of_week,
		tag
	ORDER BY
		day_of_week;
		`, days, days, days, days)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	args := []any{userID}
	if tagToQuery != "" {
		args = append(args, tagToQuery)
	}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	tagOccurrences := make(map[string][]KeyValue)

	for rows.Next() {
		var dayOfWeek, occurrences int
		var tag string
		err := rows.Scan(&dayOfWeek, &tag, &occurrences)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		tagOccurrences[strconv.Itoa(dayOfWeek)] = append(tagOccurrences[strconv.Itoa(dayOfWeek)], KeyValue{Key: tag, Value: occurrences})
	}

	for i := 0; i <= 6; i++ {
		if _, exists := tagOccurrences[strconv.Itoa(i)]; !exists {
			tagOccurrences[strconv.Itoa(i)] = []KeyValue{}
		}
	}

	return tagOccurrences, nil
}
