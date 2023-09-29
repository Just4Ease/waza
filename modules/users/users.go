package users

import (
	"context"
	"errors"
	"waza/models"
	"waza/repository"
)

var (
	ErrDuplicateUser       = errors.New("Sorry, user already exists")
	ErrUserNotFoundById    = errors.New("Sorry, user not found by id")
	ErrUserNotFoundByEmail = errors.New("Sorry, user not found by email")
	ErrUserNotFoundByPhone = errors.New("Sorry, user not found by phone")
)

type UserModule struct {
	userRepository repository.UserRepository
	//logger         log.Logger
}

func NewUserModule() *UserModule {
	return &UserModule{}
}

func (u UserModule) CreateUser(ctx context.Context, payload models.User) (*models.User, error) {
	// TODO: Validation on user module.

	user, err := u.userRepository.CreateUser(ctx, payload)
	if err != nil {
		// TODO: Add logger here...
		if errors.Is(err, repository.ErrDuplicateFound) {
			return nil, ErrDuplicateUser
		}
	}

	return user, nil
}

func (u UserModule) GetUserById(ctx context.Context, id string) (*models.User, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		// TODO: Add logger here...
		return nil, ErrUserNotFoundById
	}

	return user, nil
}

func (u UserModule) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		// TODO: Add logger here...
		return nil, ErrUserNotFoundByEmail
	}

	return user, nil
}

func (u UserModule) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	user, err := u.userRepository.GetUserByPhone(ctx, phone)
	if err != nil {
		// TODO: Add logger here...
		return nil, ErrUserNotFoundByPhone
	}

	return user, nil
}
