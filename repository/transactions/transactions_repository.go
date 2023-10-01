package transactionsrepository

import (
	"context"
	"database/sql"
	"waza/models"
	"waza/repository"
)

const tableSetup = `
CREATE TABLE IF NOT EXISTS transactions (
    	id TEXT PRIMARY KEY UNIQUE,
		amount FLOAT,
		description TEXT,
		reference TEXT,
		sourceAccountId TEXT,
		sourceAccountName TEXT,
		destinationAccountId TEXT,
		destinationAccountName TEXT,
		status TEXT,
		type TEXT,
		balanceBeforeCredit FLOAT,
		balanceAfterCredit FLOAT,
		balanceBeforeDebit FLOAT,
		balanceAfterDebit FLOAT,
		timeCreated DATETIME,
		timeUpdated DATETIME
)
`

type transactionsRepo struct {
	dataStore *sql.DB
}

func (t transactionsRepo) CreateTransaction(ctx context.Context, payload *models.Transaction) (*models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (t transactionsRepo) GetTransactionById(ctx context.Context, id string) (*models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (t transactionsRepo) UpdateTransaction(ctx context.Context, payload *models.Transaction) (*models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (t transactionsRepo) ListTransactionHistoryByAccountId(ctx context.Context, accountId string) ([]*models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func NewTransactionsRepository(db *sql.DB) (repository.TransactionRepository, error) {
	if _, err := db.Exec(tableSetup); err != nil {
		return nil, err
	}

	return &transactionsRepo{
		dataStore: db,
	}, nil
}
