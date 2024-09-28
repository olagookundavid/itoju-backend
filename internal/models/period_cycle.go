package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type MenstrualCycle struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	StartDate    time.Time `json:"start_date"`
	CycleLength  int       `json:"cycle_length"`
	PeriodLength int       `json:"period_length"`
}

type CycleDay struct {
	ID          string    `json:"id"`
	CycleID     string    `json:"cycle_id"`
	UserID      string    `json:"-"`
	Date        time.Time `json:"date"`
	IsPeriod    bool      `json:"is_period"`
	IsOvulation bool      `json:"is_ovulation"`
	Flow        float32   `json:"flow"`
	Pain        float32   `json:"pain"`
	Tags        []string  `json:"tags"`
	CMQ         string    `json:"cmq"`
}

type UserPeriodModel struct {
	DB *sql.DB
}

func (m *UserPeriodModel) ReturnMenstrualCycle(userID string, cycleLength, periodLength int, startDate time.Time) MenstrualCycle {
	return MenstrualCycle{
		UserID: userID, StartDate: startDate, CycleLength: cycleLength, PeriodLength: periodLength,
	}
}

func (m *UserPeriodModel) ReturnCycleDay(cycleID, userID string, isPeriod, isOvulation bool, date time.Time) CycleDay {
	return CycleDay{
		CycleID: cycleID, UserID: userID, IsPeriod: isPeriod, IsOvulation: isOvulation, Date: date, CMQ: "", Tags: []string{},
	}
}

// func (m *UserPeriodModel) ReturnCycleDay(CMQ string, cycleID, flow, pain int, isPeriod, isOvulation bool, date time.Time, tags []string) CycleDay {
// 	return CycleDay{
// 		CMQ: CMQ, CycleID: cycleID, Flow: flow, Pain: pain, IsPeriod: isPeriod, IsOvulation: isOvulation, Date: date,
// 	}
// }

func (m *UserPeriodModel) GetMenstrualCycles(userID string) ([]MenstrualCycle, error) {
	query := `SELECT id, user_id, start_date, cycle_length, period_length
              FROM menstrual_cycles WHERE user_id = $1 ORDER BY start_date DESC`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cycles []MenstrualCycle
	for rows.Next() {
		var cycle MenstrualCycle
		err := rows.Scan(&cycle.ID, &cycle.UserID, &cycle.StartDate, &cycle.CycleLength, &cycle.PeriodLength)
		if err != nil {
			return nil, err
		}
		cycles = append(cycles, cycle)
	}
	return cycles, nil
}

func (m *UserPeriodModel) GetCycleDays(cycleID, userID string) ([]CycleDay, error) {

	if cycleID == "" {
		return []CycleDay{}, nil
	}

	query := `SELECT id, cycle_id, date, is_period, is_ovulation, flow, pain, tags, cmq
              FROM cycles_days WHERE cycle_id = $1 AND user_id = $2 ORDER BY date ASC`

	rows, err := m.DB.Query(query, cycleID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []CycleDay
	for rows.Next() {
		var day CycleDay
		err := rows.Scan(&day.ID, &day.CycleID, &day.Date, &day.IsPeriod, &day.IsOvulation, &day.Flow, &day.Pain, pq.Array(&day.Tags), &day.CMQ)
		if err != nil {
			return nil, err
		}
		days = append(days, day)
	}
	if days == nil {
		days = []CycleDay{}
	}
	return days, nil
}

func (m *UserPeriodModel) InsertMenstrualCycleTx(tx *sql.Tx, cycle *MenstrualCycle) (string, error) {
	query := `INSERT INTO menstrual_cycles (user_id, start_date, cycle_length, period_length, created_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id string
	err := tx.QueryRowContext(ctx, query, cycle.UserID, cycle.StartDate, cycle.CycleLength, cycle.PeriodLength, time.Now()).Scan(&id)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "menstrual_cycles_start_date_key"`:
			return "", ErrRecordAlreadyExist
		default:
			return "", err
		}
	}
	return id, nil
}

func (m *UserPeriodModel) InsertCycleDayTx(tx *sql.Tx, day *CycleDay) error {
	query := `INSERT INTO cycles_days (cycle_id, user_id, date, is_period, is_ovulation, flow, pain, tags, cmq, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, day.CycleID, day.UserID, day.Date, day.IsPeriod, day.IsOvulation, day.Flow, day.Pain, pq.Array(day.Tags), day.CMQ, time.Now())
	if err != nil {

		print(err.Error())
		return err
	}
	return nil
}

func (m *UserPeriodModel) UpdateCycleDay(cycleDay *CycleDay) error {

	query := `UPDATE cycles_days SET flow = $1, pain = $2, is_ovulation = $3, is_period = $4, cmq = $5, tags = $6  WHERE id = $7`

	args := []any{cycleDay.Flow, cycleDay.Pain, cycleDay.IsOvulation, cycleDay.IsPeriod, cycleDay.CMQ, pq.Array(cycleDay.Tags), cycleDay.ID}
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

func (m *UserPeriodModel) GetCycleDay(id string) (*CycleDay, error) {
	if id == "" {
		return nil, ErrRecordNotFound
	}
	query := ` SELECT id, cycle_id, date, is_period, is_ovulation, flow, pain, tags, cmq FROM cycles_days WHERE id = $1; `
	var cycleDay CycleDay
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&cycleDay.ID,
		&cycleDay.CycleID,
		&cycleDay.Date,
		&cycleDay.IsPeriod,
		&cycleDay.IsOvulation,
		&cycleDay.Flow,
		&cycleDay.Pain,
		pq.Array(&cycleDay.Tags),
		&cycleDay.CMQ,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &cycleDay, nil
}

func (m *UserPeriodModel) GetMensesCycleIds(id string) ([]string, error) {
	if id == "" {
		return nil, ErrRecordNotFound
	}
	query := ` SELECT id FROM menstrual_cycles WHERE user_id = $1 ORDER BY created_at DESC LIMIT 3 `
	var cycleDayIds []string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cycleDayId string
		if err := rows.Scan(&cycleDayId); err != nil {
			return nil, err
		}
		cycleDayIds = append(cycleDayIds, cycleDayId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cycleDayIds, nil
}
