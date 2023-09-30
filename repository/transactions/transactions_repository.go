package transactionsrepository

import (
	"github.com/tidwall/buntdb"
	"waza/repository"
)

type transactionsRepo struct {
	dataStore *buntdb.DB
}

func NewTransactionsRepository(dataStorageFile string) (repository.TransactionRepository, error) {
	db, err := buntdb.Open(dataStorageFile)
	if err != nil {
		return nil, err
	}

	return &transactionsRepo{
		dataStore: db,
	}, nil
}
