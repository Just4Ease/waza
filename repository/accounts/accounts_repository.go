package accountsrepository

import (
	"context"
	"github.com/tidwall/buntdb"
	"path"
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

func NewAccountRepository(dataStorageDir string) (repository.AccountRepository, error) {
	database := path.Join(dataStorageDir, "accounts.db")

	db, err := buntdb.Open(database)
	if err != nil {
		return nil, err
	}

	return &accountsRepo{
		dataStore: db,
	}, nil
}
