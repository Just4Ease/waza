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
}

type TransactionRepository interface {
	//Credit()
}
