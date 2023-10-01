package transactionsrepository

import (
	"context"
	"database/sql"
	"time"
	"waza/models"
	"waza/repository"
	"waza/utils"
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

func (t transactionsRepo) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	transaction.Id = utils.GenerateId()
	now := time.Now()

	transaction.TimeCreated = now
	transaction.TimeUpdated = now

	insertSQL := `
        INSERT INTO transactions (id, amount, description, reference, sourceAccountId,
            sourceAccountName, destinationAccountId, destinationAccountName, status,
            type, balanceBeforeCredit, balanceAfterCredit, balanceBeforeDebit, balanceAfterDebit,
            timeCreated, timeUpdated)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := t.dataStore.ExecContext(ctx, insertSQL, transaction.Id, transaction.Amount, transaction.Description,
		transaction.Reference, transaction.SourceAccountId, transaction.SourceAccountName,
		transaction.DestinationAccountId, transaction.DestinationAccountName, transaction.Status,
		transaction.Type, transaction.BalanceBeforeCredit, transaction.BalanceAfterCredit,
		transaction.BalanceBeforeDebit, transaction.BalanceAfterDebit, transaction.TimeCreated,
		transaction.TimeUpdated)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t transactionsRepo) GetTransactionById(ctx context.Context, id string) (*models.Transaction, error) {
	row := t.dataStore.QueryRowContext(ctx, "SELECT * from transactions WHERE id = ?", id)
	return scanner(row)
}

func (t transactionsRepo) UpdateTransaction(ctx context.Context, transactionId string, status models.TransactionStatus) (*models.Transaction, error) {
	now := time.Now()
	// Define the SQL query for updating a transaction by ID.
	updateSQL := `
        UPDATE transactions
        SET timeUpdated = ?, status = ?
        WHERE id = ?`

	_, err := t.dataStore.ExecContext(ctx, updateSQL, now, status, transactionId)
	if err != nil {
		return nil, err
	}
	return t.GetTransactionById(ctx, transactionId)
}

func (t transactionsRepo) ListTransactionHistoryByAccountId(ctx context.Context, accountId string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	// Define the SQL query for listing transaction history by source account ID.
	querySQL := `
        SELECT * FROM transactions WHERE sourceAccountId = ?`

	rows, err := t.dataStore.QueryContext(ctx, querySQL, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		transaction, err := scanRows(rows)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func NewTransactionsRepository(db *sql.DB) (repository.TransactionRepository, error) {
	if _, err := db.Exec(tableSetup); err != nil {
		return nil, err
	}

	return &transactionsRepo{
		dataStore: db,
	}, nil
}

func scanner(row *sql.Row) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := row.Scan(
		&transaction.Id,
		&transaction.Amount,
		&transaction.Description,
		&transaction.Reference,
		&transaction.SourceAccountId,
		&transaction.SourceAccountName,
		&transaction.DestinationAccountId,
		&transaction.DestinationAccountName,
		&transaction.Status,
		&transaction.Type,
		&transaction.BalanceBeforeCredit,
		&transaction.BalanceAfterCredit,
		&transaction.BalanceBeforeDebit,
		&transaction.BalanceAfterDebit,
		&transaction.TimeCreated,
		&transaction.TimeUpdated,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func scanRows(rows *sql.Rows) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := rows.Scan(
		&transaction.Id,
		&transaction.Amount,
		&transaction.Description,
		&transaction.Reference,
		&transaction.SourceAccountId,
		&transaction.SourceAccountName,
		&transaction.DestinationAccountId,
		&transaction.DestinationAccountName,
		&transaction.Status,
		&transaction.Type,
		&transaction.BalanceBeforeCredit,
		&transaction.BalanceAfterCredit,
		&transaction.BalanceBeforeDebit,
		&transaction.BalanceAfterDebit,
		&transaction.TimeCreated,
		&transaction.TimeUpdated,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}
