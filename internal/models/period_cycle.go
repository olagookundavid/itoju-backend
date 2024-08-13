package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type MenstrualCycle struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	StartDate    time.Time `json:"start_date"`
	CycleLength  int       `json:"cycle_length"`
	PeriodLength int       `json:"period_length"`
}

type CycleDay struct {
	ID          int       `json:"id"`
	CycleID     int       `json:"cycle_id"`
	Date        time.Time `json:"date"`
	IsPeriod    bool      `json:"is_period"`
	IsOvulation bool      `json:"is_ovulation"`
	Flow        int       `json:"flow"`
	Pain        int       `json:"pain"`
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

func (m *UserPeriodModel) ReturnCycleDay(cycleID int, isPeriod, isOvulation bool, date time.Time) CycleDay {
	return CycleDay{
		CycleID: cycleID, IsPeriod: isPeriod, IsOvulation: isOvulation, Date: date,
	}
}

// func (m *UserPeriodModel) ReturnCycleDay(CMQ string, cycleID, flow, pain int, isPeriod, isOvulation bool, date time.Time, tags []string) CycleDay {
// 	return CycleDay{
// 		CMQ: CMQ, CycleID: cycleID, Flow: flow, Pain: pain, IsPeriod: isPeriod, IsOvulation: isOvulation, Date: date,
// 	}
// }

func (m *UserPeriodModel) InsertMenstrualCycle(cycle *MenstrualCycle) (int, error) {
	query := `INSERT INTO menstrual_cycles (user_id, start_date, cycle_length, period_length)
              VALUES ($1, $2, $3, $4) RETURNING id`

	var id int
	err := m.DB.QueryRow(query, cycle.UserID, cycle.StartDate, cycle.CycleLength, cycle.PeriodLength).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *UserPeriodModel) InsertCycleDay(day *CycleDay) (int, error) {
	query := `INSERT INTO cycles_days (cycle_id, date, is_period, is_ovulation, flow, pain, tags, cmq)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var id int
	err := m.DB.QueryRow(query, day.CycleID, day.Date, day.IsPeriod, day.IsOvulation, day.Flow, day.Pain, pq.Array(day.Tags), day.CMQ).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

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

func (m *UserPeriodModel) GetCycleDays(cycleID int) ([]CycleDay, error) {
	query := `SELECT id, cycle_id, date, is_period, is_ovulation, flow, pain, tags, cmq
              FROM cycles_days WHERE cycle_id = $1 ORDER BY date ASC`

	rows, err := m.DB.Query(query, cycleID)
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
	return days, nil
}

func (m *UserPeriodModel) InsertMenstrualCycleTx(tx *sql.Tx, cycle *MenstrualCycle) (int, error) {
	query := `INSERT INTO menstrual_cycles (user_id, start_date, cycle_length, period_length, created_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	err := tx.QueryRowContext(ctx, query, cycle.UserID, cycle.StartDate, cycle.CycleLength, cycle.PeriodLength, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *UserPeriodModel) InsertCycleDayTx(tx *sql.Tx, day *CycleDay) error {
	query := `INSERT INTO cycles_days (cycle_id, date, is_period, is_ovulation, flow, pain, tags, cmq, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, day.CycleID, day.Date, day.IsPeriod, day.IsOvulation, day.Flow, day.Pain, pq.Array(day.Tags), day.CMQ, time.Now())
	if err != nil {
		return err
	}
	return nil
}
