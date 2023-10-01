package transactions

import (
	"github.com/sirupsen/logrus"
	"waza/repository"
	"waza/store"
)

type TransactionService struct {
	eventStore            store.EventStore
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
	logger                *logrus.Logger
}

func NewTransactionService(
	eventStore store.EventStore,
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
	logger *logrus.Logger,
) *TransactionService {
	return &TransactionService{
		eventStore:            eventStore,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		logger:                logger,
	}
}
