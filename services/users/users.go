package users

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"waza/models"
	"waza/repository"
)

var (
	ErrDuplicateUser       = errors.New("Sorry, user already exists")
	ErrUserNotFoundById    = errors.New("Sorry, user not found by id")
	ErrUserNotFoundByEmail = errors.New("Sorry, user not found by email")
	ErrUserNotFoundByPhone = errors.New("Sorry, user not found by phone")
)

type UserService struct {
	userRepository repository.UserRepository
	logger         *logrus.Logger
}

func NewUserService(userRepository repository.UserRepository, logger *logrus.Logger) *UserService {
	return &UserService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (u UserService) CreateUser(ctx context.Context, payload models.User) (*models.User, error) {
	// TODO: Validation on user module.

	user, err := u.userRepository.CreateUser(ctx, payload)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to create user")
		if errors.Is(err, repository.ErrDuplicateFound) {
			return nil, ErrDuplicateUser
		}
	}

	return user, nil
}

func (u UserService) GetUserById(ctx context.Context, id string) (*models.User, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to get user by id")
		return nil, ErrUserNotFoundById
	}

	return user, nil
}

func (u UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to get user by email")
		return nil, ErrUserNotFoundByEmail
	}

	return user, nil
}

func (u UserService) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	user, err := u.userRepository.GetUserByPhone(ctx, phone)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to get user by phone")
		return nil, ErrUserNotFoundByPhone
	}

	return user, nil
}
