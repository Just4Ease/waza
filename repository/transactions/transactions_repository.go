package transactionsrepository

import (
	"database/sql"
	"waza/repository"
)

type transactionsRepo struct {
	dataStore *sql.DB
}

func NewTransactionsRepository(db *sql.DB) (repository.TransactionRepository, error) {

	return &transactionsRepo{
		dataStore: db,
	}, nil
}
