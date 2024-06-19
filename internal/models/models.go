package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrRecordAlreadyExist = errors.New("already exists")
)

type Models struct {
	Users            UserModel
	Tokens           TokenModel
	Metrics          MetricsModel
	Smileys          SmileysModel
	Symptoms         SymptomsModel
	Conditions       ConditionsModel
	Resources        ResourcesModel
	Menses           MensesModels
	BodyMeasure      BodyMeasureModel
	SymsMetric       SymsMetricModel
	SleepMetric      SleepMetricModel
	FoodMetric       FoodMetricModel
	ExerciseMetric   ExerciseMetricModel
	UrineMetric      UrineMetricModel
	BowelMetric      BowelMetricModel
	MedicationMetric MedicationMetricModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:            UserModel{DB: db},
		Tokens:           TokenModel{DB: db},
		Metrics:          MetricsModel{DB: db},
		Smileys:          SmileysModel{DB: db},
		Symptoms:         SymptomsModel{DB: db},
		Conditions:       ConditionsModel{DB: db},
		Resources:        ResourcesModel{DB: db},
		Menses:           MensesModels{DB: db},
		BodyMeasure:      BodyMeasureModel{DB: db},
		SymsMetric:       SymsMetricModel{DB: db},
		SleepMetric:      SleepMetricModel{DB: db},
		FoodMetric:       FoodMetricModel{DB: db},
		ExerciseMetric:   ExerciseMetricModel{DB: db},
		UrineMetric:      UrineMetricModel{DB: db},
		BowelMetric:      BowelMetricModel{DB: db},
		MedicationMetric: MedicationMetricModel{DB: db},
	}
}
