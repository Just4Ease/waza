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
		// Now, we have to return the money to the source.
		// Hence, the reason for debit first, because we can control a reversal from here.
		_, _ = t.processDebitReversal(ctx, debitTransaction)
		// Mark the debitedTransaction as failed.

		debitTransaction.Status = models.Failed
		_, err := t.transactionRepository.UpdateTransaction(ctx, debitTransaction.Id, debitTransaction.Status)
		if err != nil {
			t.logger.WithContext(ctx).WithError(err).Error("failed to mark debited transaction as failed")
			// TODO: Report to Waza's slack, email, intercom or any channel...
			return debitTransaction, nil
		}

		return debitTransaction, nil
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

	// Deliberately Wrote this update here...
	// Pheww! submitting test )
	debitTransaction.Status = models.Completed
	completedDebitTransaction, err := t.transactionRepository.UpdateTransaction(ctx, debitTransaction.Id, debitTransaction.Status)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("failed to mark debited transaction as failed")
		// TODO: Report to Waza's slack, email, intercom or any channel...
		return debitTransaction, nil
	}
	t.publishCompletedTransaction(ctx, completedDebitTransaction)

	// We are returning the debit transaction.
	// because this is the source more like when a person initiates a transfer, they want to see their receipt.
	return completedDebitTransaction, nil
}

// processDebitReversal is an internal function that can be trusted to reverse a debit on an account.
// Once a debit has been reversed, a new transaction of status `reversed` is created, and the previous transaction of status `pending` is marked as `failed`
func (t TransactionService) processDebitReversal(ctx context.Context, debitTransaction *models.Transaction) (*models.Transaction, error) {
	account, err := t.accountRepository.Credit(ctx, debitTransaction.SourceAccountId, debitTransaction.Amount)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("failed to reverse funds back to the debited account.")
		// TODO: Push this to an event queue for later processing or report to Waza finance admin via slack, email, and intercom.
		return nil, err
	}

	// We are creating a completed credit transaction, because that's the only reasonable status to have.
	now := time.Now()
	reversedTransaction, err := t.transactionRepository.CreateTransaction(ctx, &models.Transaction{
		TimeCreated:            now,
		TimeUpdated:            now,
		Reference:              debitTransaction.Reference,
		Description:            fmt.Sprintf("REVERSAL %s", debitTransaction.Description),
		Amount:                 debitTransaction.Amount,
		Type:                   models.Credit,
		Status:                 models.Reversed,
		SourceAccountId:        "GENERAL_LEDGER",
		SourceAccountName:      "WAZA GENERAL LEDGER",
		DestinationAccountId:   debitTransaction.SourceAccountId,
		DestinationAccountName: debitTransaction.SourceAccountName,
		BalanceBeforeCredit:    account.BalanceBeforeCredit, // We can maintain the snapshot balance from the pending transaction or actual account balance snapshot, in case new money entered.
		BalanceAfterCredit:     account.BalanceAfterCredit,  // // We can maintain the snapshot balance from the pending transaction or actual account balance snapshot, in case new money entered.
		BalanceBeforeDebit:     0,
		BalanceAfterDebit:      0,
	})
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("failed to create reversed transaction record after reversing funds back to the debited account.")
		// TODO: Push this to an event queue for later processing or report to Waza finance admin via slack, email, and intercom.
		return nil, err
	}

	t.publishReversedTransaction(ctx, reversedTransaction)
	return reversedTransaction, nil
}

func (t TransactionService) ListTransactionHistory(ctx context.Context, accountId string) ([]*models.Transaction, error) {
	return t.transactionRepository.ListTransactionHistoryByAccountId(ctx, accountId)
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
