package accountsrepository

import (
	"context"
	"github.com/tidwall/buntdb"
	"waza/models"
	"waza/repository"
)

type accountsRepo struct {
	dataStore *buntdb.DB
}

func (a accountsRepo) CreateAccount(ctx context.Context, payload models.Account) (*models.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a accountsRepo) GetUserById(ctx context.Context, id string) (*models.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a accountsRepo) ListAccountsByOwnerId(ctx context.Context, ownerId string) ([]*models.Account, error) {
	//TODO implement me
	panic("implement me")
}

func NewAccountRepository(dataStorageFile string) (repository.AccountRepository, error) {
	db, err := buntdb.Open(dataStorageFile)
	if err != nil {
		return nil, err
	}

	return &accountsRepo{
		dataStore: db,
	}, nil
}
