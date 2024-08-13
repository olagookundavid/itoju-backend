package models

import "database/sql"

type TransactionModel struct {
	DB *sql.DB
}

func (m *TransactionModel) BeginTx() (*sql.Tx, error) {
	return m.DB.Begin()
}
