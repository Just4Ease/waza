package accounts

import (
	"github.com/sirupsen/logrus"
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
