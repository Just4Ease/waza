package usersrepository

import (
	"context"
	"github.com/tidwall/buntdb"
	"waza/models"
	"waza/repository"
)

const usersFile = "users.db"

type userRepo struct {
	dataStore *buntdb.DB
}

func (u *userRepo) CreateUser(ctx context.Context, payload models.User) (*models.User, error) {
	panic("implement me")
}

func (u *userRepo) GetUserById(ctx context.Context, id string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	//var user models.User
	panic("implement me")
}

func (u *userRepo) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(dataStorageFile string) (repository.UserRepository, error) {
	db, err := buntdb.Open(dataStorageFile)
	if err != nil {
		return nil, err
	}

	return &userRepo{
		dataStore: db,
	}, nil
}

func buildKey(id string, phone string, email string) string {
	// id_phone_email
	return id + "_" + phone + "_" + email // I used + instead of fmt.Sprintf() for
}
