package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserPointModel struct {
	DB *sql.DB
}

func (m UserPointModel) GetUserTotalPoint(userId string, sendResult chan<- int) {
	query := ` SELECT tot_point FROM user_point 
	WHERE user_id = $1`
	var userPoint *int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, userId).Scan(
		&userPoint)

	if err != nil {
		*userPoint = 0
		print(err.Error())

	}

	sendResult <- *userPoint
}

func (m UserPointModel) GetUserTotalPoints(userId string, sendDayResult chan<- int, sendMonthResult chan<- int) {
	query := `
	SELECT 
		COALESCE(SUM(point) FILTER (WHERE date_trunc('day', date) = CURRENT_DATE), 0) AS today_points, 
		COALESCE(SUM(point) FILTER (WHERE date_trunc('week', date) = date_trunc('week', CURRENT_DATE)), 0) AS this_week_points
	FROM user_point_record
	WHERE user_id = $1 `

	var userDayPoint, userMonthPoint *int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, userId).Scan(
		&userDayPoint,
		&userMonthPoint)

	if err != nil {
		*userDayPoint = 0
		*userMonthPoint = 0

	}
	sendDayResult <- *userDayPoint
	sendMonthResult <- *userMonthPoint
}

func (m UserPointModel) InsertPoint(userId, scope string, point int64) error {

	tx, err := m.DB.Begin()
	if err != nil {
		return errors.New("could add user point")
	}

	// Ensure to rollback the transaction in case of an error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	query := `
	INSERT INTO user_point (user_id, tot_point)
	VALUES ($1, $2)
	ON CONFLICT (user_id)
	DO UPDATE SET
    tot_point = user_point.tot_point + EXCLUDED.tot_point;
`
	userRecordQuery := `
	INSERT INTO user_point_record (user_id, point, scope)
	VALUES ($1, $2, $3);
`

	args := []any{userId, point}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	args = append(args, scope)
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = tx.ExecContext(ctx, userRecordQuery, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("couldn't add user point")
	}
	return nil
}

func (m UserPointModel) DeletePointRecordMoreThanWeek() error {
	query := `DELETE FROM user_point_record WHERE date < CURRENT_DATE - INTERVAL '7 days'`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query)
	return err
}
