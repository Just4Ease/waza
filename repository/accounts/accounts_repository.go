package accountsrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
	"waza/models"
	"waza/repository"
	"waza/utils"
)

const tableSetup = `
CREATE TABLE IF NOT EXISTS accounts (
    	id TEXT PRIMARY KEY UNIQUE,
		accountName TEXT,
		accountOwnerId TEXT,
		currency TEXT,
		iso2 TEXT,
		balance FLOAT,
		balanceBeforeCredit FLOAT,
		balanceAfterCredit FLOAT,
		balanceBeforeDebit FLOAT,
		balanceAfterDebit FLOAT,
		timeCreated DATETIME,
		timeUpdated DATETIME
)
`

type accountsRepo struct {
	dataStore *sql.DB
}

func (a accountsRepo) CreateAccount(ctx context.Context, payload *models.Account) (*models.Account, error) {
	// NOTE: I did not bother checking for duplicate here,
	// as the only account creation method is when a non-duplicate user is successfully created.
	now := time.Now()

	payload.Id = utils.GenerateId()
	payload.TimeCreated = now
	payload.TimeUpdated = now

	_txt := `
INSERT INTO 
    accounts (id, accountName, accountOwnerId, currency, iso2, balance, balanceBeforeCredit, balanceAfterCredit, balanceBeforeDebit, balanceAfterDebit,  timeCreated, timeUpdated)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	statement, err := a.dataStore.PrepareContext(ctx, _txt)
	if err != nil {
		return nil, err
	}

	_, err = statement.ExecContext(ctx,
		payload.Id,
		payload.AccountName,
		payload.AccountOwnerId,
		payload.Currency,
		payload.Iso2,
		payload.Balance,
		payload.TimeCreated,
		payload.TimeUpdated,
	)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a accountsRepo) GetAccountById(ctx context.Context, id string) (*models.Account, error) {
	row := a.dataStore.QueryRowContext(ctx, "SELECT * from accounts WHERE id = ?", id)
	return scanner(row)
}

func (a accountsRepo) GetAccountByOwnerId(ctx context.Context, ownerId string) (*models.Account, error) {
	fmt.Print(ownerId)
	row := a.dataStore.QueryRowContext(ctx, "SELECT * from accounts WHERE accountOwnerId = ?", ownerId)
	return scanner(row)
}

func (a accountsRepo) Credit(ctx context.Context, accountId string, amount float64) (*models.Account, error) {
	if amount <= 0 || amount >= math.MaxFloat64 {
		return nil, errors.New("amount has to be within a valid range")
	}

	tx, err := a.dataStore.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelLinearizable,
		ReadOnly:  false,
	})
	if err != nil {
		log.Println("Error beginning transaction in Credit Account method.", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback the transaction in case of a panic.
			if err := tx.Rollback(); err != nil {
				log.Println("Error returned from transaction rollback in Credit Account method.", err)
			}
		}
	}()

	row := tx.QueryRowContext(ctx, "SELECT * from accounts WHERE id = ?", accountId)
	account, err := scanner(row)
	if err != nil {
		return nil, err
	}

	newBalance := account.Balance + amount
	account.BalanceBeforeCredit = account.Balance
	account.BalanceAfterCredit = newBalance
	account.Balance = newBalance
	account.TimeUpdated = time.Now()

	if err = tx.Commit(); err != nil {
		log.Println("Error committing credit transaction:", err)
		if err := tx.Rollback(); err != nil {
			log.Println("Error returned from transaction rollback in Credit Account method when transaction failed to commit.", err)
			return nil, err
		}
		return nil, err
	}

	// This account here will be used to record or update a transaction log.
	return account, nil
}

func (a accountsRepo) Debit(ctx context.Context, accountId string, amount float64) (*models.Account, error) {
	if amount <= 0 || amount >= math.MaxFloat64 {
		return nil, errors.New("amount has to be within a valid range")
	}

	tx, err := a.dataStore.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelLinearizable,
		ReadOnly:  false,
	})
	if err != nil {
		log.Println("Error beginning transaction in Debit Account method.", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback the transaction in case of a panic.
			if err := tx.Rollback(); err != nil {
				log.Println("Error returned from transaction rollback in Debit Account method.", err)
			}
		}
	}()

	row := tx.QueryRowContext(ctx, "SELECT * from accounts WHERE id = ?", accountId)
	account, err := scanner(row)
	if err != nil {
		return nil, err
	}

	if amount > account.Balance {
		return nil, errors.New("insufficient balance")
	}

	account.BalanceBeforeDebit = account.Balance
	account.Balance -= amount
	account.BalanceAfterDebit = account.Balance
	account.TimeUpdated = time.Now()

	if err = tx.Commit(); err != nil {
		log.Println("Error committing debit account transaction:", err)
		if err := tx.Rollback(); err != nil {
			log.Println("Error returned from transaction rollback in Debit Account method when transaction failed to commit.", err)
			return nil, err
		}
		return nil, err
	}

	// This account here will be used to record or update a transaction log.
	return account, nil
}

func NewAccountRepository(db *sql.DB) (repository.AccountRepository, error) {
	// Setup table, if not exist.
	if _, err := db.Exec(tableSetup); err != nil {
		return nil, err
	}

	return &accountsRepo{
		dataStore: db,
	}, nil
}

func scanner(row *sql.Row) (*models.Account, error) {
	var account models.Account
	if err := row.Scan(
		&account.Id,
		&account.AccountName,
		&account.AccountOwnerId,
		&account.Currency,
		&account.Iso2,
		&account.Balance,
		&account.BalanceBeforeCredit,
		&account.BalanceAfterCredit,
		&account.BalanceBeforeDebit,
		&account.BalanceAfterDebit,
		&account.TimeCreated,
		&account.TimeUpdated,
	); err != nil {
		return nil, err
	}

	return &account, nil
}
