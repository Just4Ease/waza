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
	Id                 string
	AccountName        string
	AccountOwnerId     string
	Currency           string
	Iso2               string
	Balance            float64
	BalanceAfterDebit  float64
	BalanceAfterCredit float64
}

type Transaction struct {
	SourceAccountId        string
	SourceAccountName      string
	DestinationAccountId   string
	DestinationAccountName string
}
