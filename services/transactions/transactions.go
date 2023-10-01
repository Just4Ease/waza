package transactions

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
	"waza/models"
	"waza/repository"
	"waza/store"
	"waza/utils"
)

var (
	ErrFetchingSourceAccount      = errors.New("Sorry, failed to get source account for transfer")
	ErrFetchingDestinationAccount = errors.New("Sorry, failed to get destination account for transfer")
)

type TransactionService struct {
	eventStore            store.EventStore
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
	logger                *logrus.Logger
}

func NewTransactionService(
	eventStore store.EventStore,
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
	logger *logrus.Logger,
) *TransactionService {
	return &TransactionService{
		eventStore:            eventStore,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		logger:                logger,
	}
}

func (t TransactionService) TransferFunds(ctx context.Context, payload models.TransactionPayload) (*models.Transaction, error) {
	if err := payload.Validate(); err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("funds transfer payload validation failed")
		return nil, err
	}

	// Validate account existence...
	sourceAccount, err := t.accountRepository.GetAccountById(ctx, payload.FromAccountId)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("failed to fetch source account for transfer.")
		return nil, ErrFetchingSourceAccount
	}

	destinationAccount, err := t.accountRepository.GetAccountById(ctx, payload.ToAccountId)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("failed to fetch destination account for transfer.")
		return nil, ErrFetchingDestinationAccount
	}

	// I'm following the simple rule:
	// Debit source first.
	// Credit destination second.
	// If destination Credit fails, Reverse Debit on Source.

	now := time.Now()
	reference := generateTransactionReference(now)

	debitedAccount, err := t.accountRepository.Debit(ctx, sourceAccount.Id, payload.Amount)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("failed to debit account")
		return nil, err
	}

	debitTransaction, err := t.transactionRepository.CreateTransaction(ctx, &models.Transaction{
		TimeCreated:            now,
		TimeUpdated:            now,
		Reference:              reference,
		Description:            payload.Description,
		Amount:                 payload.Amount,
		Type:                   models.Debit,
		Status:                 models.Pending,
		SourceAccountId:        sourceAccount.Id,
		SourceAccountName:      sourceAccount.AccountName,
		DestinationAccountId:   destinationAccount.Id,
		DestinationAccountName: destinationAccount.AccountName,
		BalanceBeforeDebit:     debitedAccount.BalanceBeforeDebit,
		BalanceAfterDebit:      debitedAccount.BalanceAfterDebit,
		BalanceAfterCredit:     0,
		BalanceBeforeCredit:    0,
	})
	// Ideally, I would publish this transaction created event to the event store.
	// ( For log keeping, or to trigger AML watchers if amount is beyond a certain threshold and other business or government centric constraints )
	t.publishCreatedTransaction(ctx, debitTransaction)

	creditedAccount, err := t.accountRepository.Credit(ctx, destinationAccount.Id, payload.Amount)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("failed to credit destination account")
		t.logger.WithContext(ctx).Infof("rolling back debit on source account: %s", payload.FromAccountId)

		// Ok comrade... We couldn't credit the destination account. But we debited the source account.
		// Now, we have to return the money to the source. ( hence the reason for debit first, because we can control a reversal from here. )
		// Publish transaction failed event here...
		return nil, err
	}

	// We are creating a completed credit transaction, because that's the only reasonable status to have.
	creditTransaction, err := t.transactionRepository.CreateTransaction(ctx, &models.Transaction{
		TimeCreated:            now,
		TimeUpdated:            now,
		Reference:              reference,
		Description:            payload.Description,
		Amount:                 payload.Amount,
		Type:                   models.Credit,
		Status:                 models.Completed,
		SourceAccountId:        sourceAccount.Id,
		SourceAccountName:      sourceAccount.AccountName,
		DestinationAccountId:   destinationAccount.Id,
		DestinationAccountName: destinationAccount.AccountName,
		BalanceBeforeCredit:    creditedAccount.BalanceBeforeCredit,
		BalanceAfterCredit:     creditedAccount.BalanceAfterCredit,
		BalanceBeforeDebit:     0,
		BalanceAfterDebit:      0,
	})

	// We publish only the credited transaction data.
	// Then, use this to go and update the debit transaction in the background.
	t.publishCompletedTransaction(ctx, creditTransaction)

	// We are returning the debit transaction.
	// because this is the source more like when a person initiates a transfer, they want to see their receipt.
	return debitTransaction, nil
}

func (t TransactionService) ProcessDebitReversal(ctx context.Context, failedDebitTransaction *models.Transaction) (*models.Transaction, error) {
	return nil, nil
}

type CreditAdvice struct {
	Account     *models.Account
	Description string
	Reference   string
}

func (t TransactionService) CreditAccount(ctx context.Context, account *models.Account, description string) {

}

func generateTransactionReference(now time.Time) string {
	YMDHMString := fmt.Sprintf("%d%d%d%d%d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
	)
	return fmt.Sprintf("%s%s", utils.GenerateId(), YMDHMString)
}
