package models

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type User struct {
	Id          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       *string   `json:"email"`
	Phone       string    `json:"phone"`
	TimeCreated time.Time `json:"timeCreated"`
	TimeUpdated time.Time `json:"timeUpdated"` // Hmm... not really useful throughout this test 💀
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Phone, validation.Required, is.E164),
	)
}

type Account struct {
	Id                  string    `json:"id"`
	AccountName         string    `json:"accountName"`
	AccountOwnerId      string    `json:"accountOwnerId"`
	Currency            string    `json:"currency"`
	Iso2                string    `json:"iso2"`
	Balance             float64   `json:"balance"`
	BalanceBeforeCredit float64   `json:"balanceBeforeCredit"` // polled from transaction.
	BalanceBeforeDebit  float64   `json:"balanceBeforeDebit"`  // polled from transaction.
	BalanceAfterDebit   float64   `json:"balanceAfterDebit"`   // polled from transaction.
	BalanceAfterCredit  float64   `json:"balanceAfterCredit"`  // polled from transaction.
	TimeCreated         time.Time `json:"timeCreated"`
	TimeUpdated         time.Time `json:"timeUpdated"`
}

type Transaction struct {
	Id                     string            `json:"id"`
	Reference              string            `json:"reference"`
	Description            string            `json:"description"`
	Amount                 float64           `json:"amount"`
	Type                   TransactionType   `json:"type"`
	Status                 TransactionStatus `json:"status"`
	SourceAccountId        string            `json:"sourceAccountId"`
	SourceAccountName      string            `json:"sourceAccountName"`
	DestinationAccountId   string            `json:"destinationAccountId"`
	DestinationAccountName string            `json:"destinationAccountName"`
	BalanceBeforeCredit    float64           `json:"balanceBeforeCredit"`
	BalanceBeforeDebit     float64           `json:"balanceBeforeDebit"`
	BalanceAfterDebit      float64           `json:"balanceAfterDebit"`
	BalanceAfterCredit     float64           `json:"balanceAfterCredit"`
	TimeCreated            time.Time         `json:"timeCreated"`
	TimeUpdated            time.Time         `json:"timeUpdated"`
}

type TransactionType string

const (
	Credit TransactionType = "CREDIT"
	Debit  TransactionType = "DEBIT"
)

func (t TransactionType) IsValid() bool {
	switch t {
	case Credit, Debit:
		return true
	default:
		return false
	}
}

func (t TransactionType) String() string {
	return string(t)
}

type TransactionStatus string

const (
	Pending   TransactionStatus = "pending"
	Completed TransactionStatus = "completed"
	Canceled  TransactionStatus = "canceled"
	Reversed  TransactionStatus = "reversed"
)

func (t TransactionStatus) IsValid() bool {
	switch t {
	case Pending, Completed, Canceled, Reversed:
		return true
	default:
		return false
	}
}

func (t TransactionStatus) String() string {
	return string(t)
}
