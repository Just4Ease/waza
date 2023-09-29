package models

import "time"

type User struct {
	Id          string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	TimeCreated time.Time
	TimeUpdated time.Time
}

type Account struct {
	Id                  string
	AccountName         string
	AccountOwnerId      string
	Currency            string
	Iso2                string
	Balance             float64
	BalanceBeforeCredit float64 // polled from transaction.
	BalanceBeforeDebit  float64 // polled from transaction.
	BalanceAfterDebit   float64 // polled from transaction.
	BalanceAfterCredit  float64 // polled from transaction.
}

type Transaction struct {
	Id                     string
	Reference              string
	Description            string
	Amount                 float64
	Type                   TransactionType
	Status                 TransactionStatus
	SourceAccountId        string
	SourceAccountName      string
	DestinationAccountId   string
	DestinationAccountName string
	BalanceBeforeCredit    float64
	BalanceBeforeDebit     float64
	BalanceAfterDebit      float64
	BalanceAfterCredit     float64
	TimeCreated            time.Time
	TimeUpdated            time.Time
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
