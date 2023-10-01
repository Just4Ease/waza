package repository

import (
	"context"
	"errors"
	"waza/models"
)

var ErrDuplicateFound = errors.New("error: duplicate found")

type UserRepository interface {
	CreateUser(ctx context.Context, payload *models.User) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, payload *models.Account) (*models.Account, error)
	GetAccountById(ctx context.Context, id string) (*models.Account, error)
	GetAccountByOwnerId(ctx context.Context, ownerId string) (*models.Account, error)
	Credit(ctx context.Context, accountId string, amount float64) (*models.Account, error)
	Debit(ctx context.Context, accountId string, amount float64) (*models.Account, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, payload *models.Transaction) (*models.Transaction, error)
	GetTransactionById(ctx context.Context, id string) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, transactionId string, status models.TransactionStatus) (*models.Transaction, error)
	ListTransactionHistoryByAccountId(ctx context.Context, accountId string) ([]*models.Transaction, error)
}
