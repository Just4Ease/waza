package transactionsrepository

import (
	"github.com/tidwall/buntdb"
	"path"
	"waza/repository"
)

type transactionsRepo struct {
	dataStore *buntdb.DB
}

func NewTransactionsRepository(dataStorageDir string) (repository.TransactionRepository, error) {
	database := path.Join(dataStorageDir, "transactions.db")

	db, err := buntdb.Open(database)
	if err != nil {
		return nil, err
	}

	return &transactionsRepo{
		dataStore: db,
	}, nil
}
