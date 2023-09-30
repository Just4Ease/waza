// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

import (
	"time"
)

type Account struct {
	ID          string  `json:"id"`
	AccountName string  `json:"accountName"`
	Balance     float64 `json:"balance"`
	User        *User   `json:"user,omitempty"`
}

type CreateUserInput struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Phone     string  `json:"phone"`
	Email     *string `json:"email,omitempty"`
}

type Result struct {
	Success bool        `json:"success"`
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Transaction struct {
	ID        string `json:"id"`
	Reference string `json:"reference"`
}

type TransferFundsInput struct {
	FromAccountID string   `json:"fromAccountId"`
	ToAccountID   string   `json:"toAccountId"`
	Amount        *float64 `json:"amount,omitempty"`
}

type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	TimeCreated time.Time `json:"timeCreated"`
	TimeUpdated time.Time `json:"timeUpdated"`
}
