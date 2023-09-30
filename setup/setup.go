package setup

import (
	"github.com/sirupsen/logrus"
	"waza/repository"
	accountsrepository "waza/repository/accounts"
	transactionsrepository "waza/repository/transactions"
	usersrepository "waza/repository/users"
	"waza/services/accounts"
	"waza/services/transactions"
	"waza/services/users"
)

type ServiceDependencies struct {
	UserService        *users.UserService
	AccountService     *accounts.AccountService
	TransactionService *transactions.TransactionService
	Logger             *logrus.Logger
}

func ConfigureServiceDependencies(logger *logrus.Logger) *ServiceDependencies {
	var userRepository repository.UserRepository
	var accountRepository repository.AccountRepository
	var transactionsRepository repository.TransactionRepository
	var err error

	if userRepository, err = usersrepository.NewUserRepository("users.db"); err != nil {
		logrus.Fatal("failed to start service, user repository error: ", err)
	}

	if accountRepository, err = accountsrepository.NewAccountRepository("accounts.db"); err != nil {
		logrus.Fatal("failed to start service, accounts repository error: ", err)
	}

	if transactionsRepository, err = transactionsrepository.NewTransactionsRepository("transactions.db"); err != nil {
		logrus.Fatal("failed to start service, transactions repository error: ", err)
	}

	opts := &ServiceDependencies{}
	opts.Logger = logger
	opts.UserService = users.NewUserService(userRepository, logger)
	opts.AccountService = accounts.NewAccountService(accountRepository, logger)
	opts.TransactionService = transactions.NewTransactionService(accountRepository, transactionsRepository, logger)
	return opts
}
