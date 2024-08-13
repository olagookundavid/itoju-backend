package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/olagookundavid/itoju/internal/jsonlog"
)

type Symptoms struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SymptomsModel struct {
	DB *sql.DB
}

func (m SymptomsModel) GetSymptoms() ([]*Symptoms, error) {
	query := ` SELECT * FROM symptoms `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return []*Symptoms{}, err
	}
	defer rows.Close()
	symptoms := []*Symptoms{}
	for rows.Next() {
		var symptom Symptoms
		err := rows.Scan(&symptom.Id, &symptom.Name)
		if err != nil {
			return []*Symptoms{}, err
		}

		symptoms = append(symptoms, &symptom)
	}
	if err = rows.Err(); err != nil {
		return []*Symptoms{}, err
	}
	return symptoms, nil
}

func (m SymptomsModel) GetUserSymptoms(userID string) ([]*Symptoms, error) {

	query := ` SELECT symptoms.id , symptoms.name FROM symptoms
	JOIN user_symptoms ON symptoms.id = user_symptoms.symptoms_id
	WHERE user_symptoms.user_id = $1; `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	symptoms := []*Symptoms{}
	for rows.Next() {
		var symptom Symptoms
		err := rows.Scan(&symptom.Id, &symptom.Name)
		if err != nil {
			return nil, err
		}
		symptoms = append(symptoms, &symptom)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return symptoms, nil
}

func (m SymptomsModel) SetUserSymptoms(tx *sql.Tx, symID int, userID string) error {

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	query := ` INSERT INTO user_symptoms (user_id, symptoms_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, userID, symID)
	if err != nil {
		logger.PrintError(fmt.Errorf("symptom error  %s", err), nil)
		return errors.New("couldn't Add New Symptoms")
	}
	return nil

}

func (m SymptomsModel) DeleteUserSymptoms(tx *sql.Tx, userId string, symptomsID int) error {

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	query := ` DELETE FROM user_symptoms
	WHERE user_id = $1
	AND symptoms_id = $2; `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userId, symptomsID)
	if err != nil {
		logger.PrintError(fmt.Errorf("symptom error  %s", err), nil)
		return errors.New("couldn't Delete Symptoms")
	}

	return nil
}
