package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users   UserModel
	Tokens  TokenModel
	Metrics MetricsModel
	Smileys SmileysModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:   UserModel{DB: db},
		Tokens:  TokenModel{DB: db},
		Metrics: MetricsModel{DB: db},
		Smileys: SmileysModel{DB: db},
	}
}
