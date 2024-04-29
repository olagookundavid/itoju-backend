package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type FoodMetric struct {
	ID             int       `json:"id"`
	UserID         string    `json:"user_id"`
	Date           time.Time `json:"date"`
	BreakfastMeal  string    `json:"breakfast_meal"`
	LunchMeal      string    `json:"lunch_meal"`
	DinnerMeal     string    `json:"dinner_meal"`
	BreakfastExtra string    `json:"breakfast_extra"`
	LunchExtra     string    `json:"lunch_extra"`
	DinnerExtra    string    `json:"dinner_extra"`
	BreakfastFruit string    `json:"breakfast_fruit"`
	LunchFruit     string    `json:"lunch_fruit"`
	DinnerFruit    string    `json:"dinner_fruit"`
	BreakfastTags  []string  `json:"breakfast_tags"`
	LunchTags      []string  `json:"lunch_tags"`
	DinnerTags     []string  `json:"dinner_tags"`
	SnackName      string    `json:"snack_name"`
	SnackTags      []string  `json:"snack_tags"`
	GlassNo        int       `json:"glass_no"`
}

type FoodMetricModel struct {
	DB *sql.DB
}

func (m FoodMetricModel) GetUserFoodMetric(userId string, date time.Time) (*FoodMetric, error) {
	query := `
	SELECT id, user_id, date, breakfast_meal, lunch_meal, dinner_meal, breakfast_extra, lunch_extra,dinner_extra, breakfast_fruit, lunch_fruit, dinner_fruit, breakfast_tags, lunch_tags, dinner_tags, snack_name, snack_tags, glass_no
	FROM user_food_metric WHERE user_id = $1 AND date = $2;
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, userId, date)

	var foodMetric FoodMetric
	err := row.Scan(
		&foodMetric.ID,
		&foodMetric.UserID,
		&foodMetric.Date,
		&foodMetric.BreakfastMeal,
		&foodMetric.LunchMeal,
		&foodMetric.DinnerMeal,
		&foodMetric.BreakfastExtra,
		&foodMetric.LunchExtra,
		&foodMetric.DinnerExtra,
		&foodMetric.BreakfastFruit,
		&foodMetric.LunchFruit,
		&foodMetric.DinnerFruit,
		pq.Array(&foodMetric.BreakfastTags),
		pq.Array(&foodMetric.LunchTags),
		pq.Array(&foodMetric.DinnerTags),
		&foodMetric.SnackName,
		pq.Array(&foodMetric.SnackTags),
		&foodMetric.GlassNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &foodMetric, nil
}

func (m FoodMetricModel) InsertFoodMetric(foodMetric *FoodMetric) error {
	query := `
        INSERT INTO user_food_metric (user_id, date, breakfast_meal, lunch_meal, dinner_meal, 
                                       breakfast_extra, lunch_extra, dinner_extra, 
                                       breakfast_fruit, lunch_fruit, dinner_fruit, 
                                       breakfast_tags, lunch_tags, dinner_tags, 
                                       snack_name, snack_tags, glass_no)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
    `
	args := []interface{}{
		foodMetric.UserID,
		foodMetric.Date,
		foodMetric.BreakfastMeal,
		foodMetric.LunchMeal,
		foodMetric.DinnerMeal,
		foodMetric.BreakfastExtra,
		foodMetric.LunchExtra,
		foodMetric.DinnerExtra,
		foodMetric.BreakfastFruit,
		foodMetric.LunchFruit,
		foodMetric.DinnerFruit,
		pq.Array(foodMetric.BreakfastTags),
		pq.Array(foodMetric.LunchTags),
		pq.Array(foodMetric.DinnerTags),
		foodMetric.SnackName,
		pq.Array(foodMetric.SnackTags),
		foodMetric.GlassNo,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
func (m FoodMetricModel) UpdateFoodMetric(foodMetric *FoodMetric) error {
	query := `
        UPDATE user_food_metric
        SET 
            breakfast_meal = $1,
            lunch_meal = $2,
            dinner_meal = $3,
            breakfast_extra = $4,
            lunch_extra = $5,
            dinner_extra = $6,
            breakfast_fruit = $7,
            lunch_fruit = $8,
            dinner_fruit = $9,
            breakfast_tags = $10,
            lunch_tags = $11,
            dinner_tags = $12,
            snack_name = $13,
            snack_tags = $14,
            glass_no = $15
        WHERE
            id = $16 AND user_id = $17 AND date = $18
    `
	args := []interface{}{
		foodMetric.BreakfastMeal,
		foodMetric.LunchMeal,
		foodMetric.DinnerMeal,
		foodMetric.BreakfastExtra,
		foodMetric.LunchExtra,
		foodMetric.DinnerExtra,
		foodMetric.BreakfastFruit,
		foodMetric.LunchFruit,
		foodMetric.DinnerFruit,
		pq.Array(foodMetric.BreakfastTags),
		pq.Array(foodMetric.LunchTags),
		pq.Array(foodMetric.DinnerTags),
		foodMetric.SnackName,
		pq.Array(foodMetric.SnackTags),
		foodMetric.GlassNo,
		foodMetric.ID,
		foodMetric.UserID,
		foodMetric.Date,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m FoodMetricModel) CheckUserEntry(userID string, date time.Time, sendbool chan<- bool) {

	query := `
	SELECT COUNT(*) AS entry_count
	FROM user_food_metric ufm
	WHERE ufm.user_id = $1 AND ufm.date = $2
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
