package users

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"waza/events/topics"
	"waza/models"
	"waza/repository"
	"waza/store"
)

var (
	ErrDuplicateUser       = errors.New("Sorry, user already exists")
	ErrCreateUserFailed    = errors.New("Sorry, failed to create user")
	ErrUserNotFoundById    = errors.New("Sorry, user not found by id")
	ErrUserNotFoundByEmail = errors.New("Sorry, user not found by email")
	ErrUserNotFoundByPhone = errors.New("Sorry, user not found by phone")
)

type UserService struct {
	eventStore     store.EventStore
	userRepository repository.UserRepository
	logger         *logrus.Logger
}

func NewUserService(eventStore store.EventStore, userRepository repository.UserRepository, logger *logrus.Logger) *UserService {
	return &UserService{
		eventStore:     eventStore,
		userRepository: userRepository,
		logger:         logger,
	}
}

func (u UserService) CreateUser(ctx context.Context, payload models.User) (*models.User, error) {
	if err := payload.Validate(); err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to validate user data before account creation")
		return nil, err
	}

	user, err := u.userRepository.CreateUser(ctx, &payload)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to create user")
		if errors.Is(err, repository.ErrDuplicateFound) {
			return nil, ErrDuplicateUser
		}

		return nil, ErrCreateUserFailed
	}

	data, _ := json.Marshal(user)
	if err := u.eventStore.Publish(topics.UserCreated, data); err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to publish user")
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
