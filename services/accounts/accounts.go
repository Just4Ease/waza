package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"waza/events/topics"
	"waza/models"
	"waza/repository"
	"waza/store"
)

var (
	ErrCreateAccountFailed      = errors.New("Sorry, could not create account at this time.")
	ErrAccountNotFoundByOwnerId = errors.New("Sorry, account not found by provided owner id")
	ErrAccountNotFoundById      = errors.New("Sorry, account not found by provided id")
)

type AccountService struct {
	eventStore        store.EventStore
	accountRepository repository.AccountRepository
	logger            *logrus.Logger
}

func NewAccountService(
	eventStore store.EventStore,
	accountRepository repository.AccountRepository,
	logger *logrus.Logger,
) *AccountService {
	return &AccountService{
		eventStore:        eventStore,
		accountRepository: accountRepository,
		logger:            logger,
	}
}

// CreateAccount automatically creates a transactional account based on the newly created user's profile.
func (a AccountService) CreateAccount(ctx context.Context, user models.User) (*models.Account, error) {
	// My intention is to make it a multi-currency account. But let me not over-engineer things :D
	account := models.Account{
		AccountName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		AccountOwnerId: user.Id,
		Currency:       "NGN", // Hard Coded :)
		Iso2:           "NG",  // Hard Coded

		// We will dash this new user's account NGN10,000 as bonus for joining Waza.
		// ( Trust Nigerians, they like free money. )
		Balance:             10000, // Hard code this for now. ( note, I'm not doing kobo related stuff for this test )
		BalanceBeforeCredit: 0,     // To be used internally for transaction
		BalanceAfterCredit:  0,     // To be used internally for transaction
		BalanceBeforeDebit:  0,     // To be used internally for transaction
		BalanceAfterDebit:   0,     // To be used internally for transaction
	}

	createdAccount, err := a.accountRepository.CreateAccount(ctx, &account)
	if err != nil {
		a.logger.WithContext(ctx).WithError(err).Error("failed to create account")
		return nil, ErrCreateAccountFailed
	}

	data, _ := json.Marshal(createdAccount)
	if err := a.eventStore.Publish(topics.AccountCreated, data); err != nil {
		a.logger.WithContext(ctx).WithError(err).Errorf("failed to publish data to %s topic.", topics.AccountCreated)
	}

	return createdAccount, nil
}

// GetAccountById returns the account by provided id
func (a AccountService) GetAccountById(ctx context.Context, id string) (*models.Account, error) {
	account, err := a.accountRepository.GetAccountById(ctx, id)
	if err != nil {
		a.logger.WithContext(ctx).WithError(err).Error("failed to get account by  id")
		return nil, ErrAccountNotFoundById
	}

	return account, nil
}

// GetAccountByOwnerId returns the account belonging to the provided owner by ownerId
func (a AccountService) GetAccountByOwnerId(ctx context.Context, ownerId string) (*models.Account, error) {
	account, err := a.accountRepository.GetAccountByOwnerId(ctx, ownerId)
	if err != nil {
		a.logger.WithContext(ctx).WithError(err).Error("failed to get account by owner id")
		return nil, ErrAccountNotFoundByOwnerId
	}

	return account, nil
}
