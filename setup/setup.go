package setup

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"waza/repository"
	accountsrepository "waza/repository/accounts"
	transactionsrepository "waza/repository/transactions"
	usersrepository "waza/repository/users"
	"waza/services/accounts"
	"waza/services/transactions"
	"waza/services/users"
	"waza/store"
)

type ServiceDependencies struct {
	EventStore         store.EventStore
	UserService        *users.UserService
	AccountService     *accounts.AccountService
	TransactionService *transactions.TransactionService
	Logger             *logrus.Logger
}

func ConfigureServiceDependencies(logger *logrus.Logger) *ServiceDependencies {
	var eventStore store.EventStore
	var userRepository repository.UserRepository
	var accountRepository repository.AccountRepository
	var transactionsRepository repository.TransactionRepository
	var err error

	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		logrus.Fatal("failed to open database: ", err)
	}

	if userRepository, err = usersrepository.NewUserRepository(db); err != nil {
		logrus.Fatal("failed to start service, user repository error: ", err)
	}

	if accountRepository, err = accountsrepository.NewAccountRepository(db); err != nil {
		logrus.Fatal("failed to start service, accounts repository error: ", err)
	}

	if transactionsRepository, err = transactionsrepository.NewTransactionsRepository(db); err != nil {
		logrus.Fatal("failed to start service, transactions repository error: ", err)
	}

	eventStore = store.NewEventStore(logger)

	opts := &ServiceDependencies{
		EventStore: eventStore,
		Logger:     logger,
	}

	opts.UserService = users.NewUserService(eventStore, userRepository, logger)
	opts.AccountService = accounts.NewAccountService(eventStore, accountRepository, logger)

	opts.TransactionService = transactions.NewTransactionService(
		eventStore,
		accountRepository,
		transactionsRepository,
		logger,
	)

	return opts
}
