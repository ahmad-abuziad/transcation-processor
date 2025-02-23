package data

import (
	"database/sql"
)

type Models struct {
	SalesTransactions SalesTransactionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		SalesTransactions: SalesTransactionModel{DB: db},
	}
}
