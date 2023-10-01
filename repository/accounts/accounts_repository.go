package accountsrepository

import (
	"context"
	"database/sql"
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
    accounts (id, accountName, accountOwnerId, currency, iso2, balance, timeCreated, timeUpdated)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

	statement, err := a.dataStore.PrepareContext(ctx, _txt)
	if err != nil {
		return nil, err
	}

	_, err = statement.ExecContext(ctx,
		payload.Id,
		payload.AccountName,
		payload.AccountOwnerId,
		payload.AccountOwnerId,
		payload.Currency,
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
	row := a.dataStore.QueryRowContext(ctx, "SELECT * from accounts WHERE accountOwnerId = ?", ownerId)
	return scanner(row)
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
		&account.TimeCreated,
		&account.TimeUpdated,
	); err != nil {
		return nil, err
	}

	return &account, nil
}
