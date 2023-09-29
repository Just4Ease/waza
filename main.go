package main

import (
	"github.com/sirupsen/logrus"
	"waza/config"
	"waza/repository"
	"waza/repository/accounts"
	"waza/repository/transactions"
	"waza/repository/users"
	"waza/services/accounts"
	"waza/services/transactions"
	"waza/services/users"
)

func main() {
	secrets := config.GetSecrets()
	logger := logrus.New()

	var userRepository repository.UserRepository
	var accountRepository repository.AccountRepository
	var transactionsRepository repository.TransactionRepository
	var err error

	if userRepository, err = usersrepository.NewUserRepository(secrets.DataStorageDir); err != nil {
		logrus.Fatal("failed to start service, user repository error: ", err)
	}

	if accountRepository, err = accountsrepository.NewAccountRepository(secrets.DataStorageDir); err != nil {
		logrus.Fatal("failed to start service, accounts repository error: ", err)
	}

	if transactionsRepository, err = transactionsrepository.NewTransactionsRepository(secrets.DataStorageDir); err != nil {
		logrus.Fatal("failed to start service, transactions repository error: ", err)
	}

	userService := users.NewUserService(userRepository, logger)
	accountService := accounts.NewAccountService(accountRepository, logger)
	transactionService := transactions.NewTransactionService(accountRepository, transactionsRepository, logger)

	// TODO: GraphQL API ( using this because of the playground )
}
