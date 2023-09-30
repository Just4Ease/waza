package accounts

import (
	"context"
	"github.com/sirupsen/logrus"
	"waza/models"
	"waza/repository"
)

type AccountService struct {
	accountRepository repository.AccountRepository
	logger            *logrus.Logger
}

func NewAccountService(accountRepository repository.AccountRepository, logger *logrus.Logger) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
		logger:            logger,
	}
}

// CreateAccount from a user's profile.
func (a AccountService) CreateAccount(ctx context.Context, user models.User) (*models.Account, error) {
	panic("implement me")
}

// ListAccounts returns a list of accounts and can be filtered by: userId, currency
func (a AccountService) ListAccounts(ctx context.Context, filters map[string]string) ([]*models.Account, error) {
	panic("implement me")
}