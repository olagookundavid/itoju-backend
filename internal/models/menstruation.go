package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Menses struct {
	Id         string `json:"-"`
	Period_len int    `json:"period_len"`
	Cycle_len  int    `json:"cycle_len"`
}

type MensesModels struct {
	DB *sql.DB
}

func (m MensesModels) GetMenses(id string) (*Menses, error) {
	if id == "" {
		return nil, ErrRecordNotFound
	}
	query := ` SELECT * FROM menstruation WHERE user_id = $1`

	var menses Menses
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&menses.Id,
		&menses.Period_len,
		&menses.Cycle_len)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &menses, nil
}

func (m MensesModels) InsertMenses(menses *Menses) error {
	query := `
	INSERT INTO menstruation (user_id, period_len, cycle_len)
	VALUES ($1, $2, $3) `

	args := []any{menses.Id, menses.Cycle_len, menses.Period_len}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m MensesModels) UpdateMenses(menses *Menses) error {

	query := ` UPDATE menstruation SET period_len = $1, cycle_len = $2 WHERE user_id = $3; `

	args := []any{menses.Period_len, menses.Cycle_len, menses.Id}
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
