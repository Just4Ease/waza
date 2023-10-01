package transactions

import (
	"context"
	"encoding/json"
	"waza/events/topics"
	"waza/models"
)

func (t TransactionService) publishCreatedTransaction(ctx context.Context, transaction *models.Transaction) {
	data, _ := json.Marshal(transaction)
	if err := t.eventStore.Publish(topics.TransactionCreated, data); err != nil {
		t.logger.WithContext(ctx).WithError(err).Errorf("failed to publish event data to %s channel", topics.TransactionCreated)
	}
}

func (t TransactionService) publishCompletedTransaction(ctx context.Context, transaction *models.Transaction) {
	data, _ := json.Marshal(transaction)
	if err := t.eventStore.Publish(topics.TransactionCompleted, data); err != nil {
		t.logger.WithContext(ctx).WithError(err).Errorf("failed to publish event data to %s channel", topics.TransactionCompleted)
	}
}
