package models

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type User struct {
	TimeCreated time.Time `json:"timeCreated"`
	TimeUpdated time.Time `json:"timeUpdated"`
	Email       *string   `json:"email"`
	Id          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Phone       string    `json:"phone"`
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
	TimeCreated         time.Time `json:"timeCreated"`
	TimeUpdated         time.Time `json:"timeUpdated"`
	Id                  string    `json:"id"`
	AccountName         string    `json:"accountName"`
	AccountOwnerId      string    `json:"accountOwnerId"`
	Currency            string    `json:"currency"`
	Iso2                string    `json:"iso2"`
	Balance             float64   `json:"balance"`
	BalanceBeforeCredit float64   `json:"balanceBeforeCredit"`
	BalanceBeforeDebit  float64   `json:"balanceBeforeDebit"`
	BalanceAfterDebit   float64   `json:"balanceAfterDebit"`
	BalanceAfterCredit  float64   `json:"balanceAfterCredit"`
}

type Transaction struct {
	TimeCreated            time.Time         `json:"timeCreated"`
	TimeUpdated            time.Time         `json:"timeUpdated"`
	DestinationAccountId   string            `json:"destinationAccountId"`
	DestinationAccountName string            `json:"destinationAccountName"`
	Type                   TransactionType   `json:"type"`
	Status                 TransactionStatus `json:"status"`
	SourceAccountId        string            `json:"sourceAccountId"`
	SourceAccountName      string            `json:"sourceAccountName"`
	Id                     string            `json:"id"`
	Reference              string            `json:"reference"`
	Description            string            `json:"description"`
	BalanceBeforeDebit     float64           `json:"balanceBeforeDebit"`
	BalanceAfterDebit      float64           `json:"balanceAfterDebit"`
	BalanceAfterCredit     float64           `json:"balanceAfterCredit"`
	BalanceBeforeCredit    float64           `json:"balanceBeforeCredit"`
	Amount                 float64           `json:"amount"`
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
