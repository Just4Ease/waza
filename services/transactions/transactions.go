package transactions

import (
	"github.com/sirupsen/logrus"
	"waza/repository"
)

type TransactionService struct {
	transactionRepository repository.TransactionRepository
	accountRepository     repository.AccountRepository
	logger                *logrus.Logger
}

func NewTransactionService(
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
	logger *logrus.Logger,
) *TransactionService {
	return &TransactionService{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		logger:                logger,
	}
}
