package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type BodyMeasure struct {
	Id     string `json:"-"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

type BodyMeasureModel struct {
	DB *sql.DB
}

func (m BodyMeasureModel) GetBodyMeasure(id string) (*BodyMeasure, error) {
	if id == "" {
		return nil, ErrRecordNotFound
	}
	query := ` SELECT * FROM bodymeasure WHERE user_id = $1`

	var bodyMeasure BodyMeasure
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&bodyMeasure.Id,
		&bodyMeasure.Height,
		&bodyMeasure.Weight)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &bodyMeasure, nil
}

func (m BodyMeasureModel) InsertBodyMeasure(bodyMeasure *BodyMeasure) error {
	query := `
	INSERT INTO bodymeasure (user_id, height, weight)
	VALUES ($1, $2, $3) `

	args := []any{bodyMeasure.Id, bodyMeasure.Height, bodyMeasure.Weight}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m BodyMeasureModel) UpdateBodyMeasure(bodyMeasure *BodyMeasure) error {

	query := ` UPDATE bodymeasure SET height = $1, weight = $2 WHERE user_id = $3; `

	args := []any{bodyMeasure.Height, bodyMeasure.Weight, bodyMeasure.Id}
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
