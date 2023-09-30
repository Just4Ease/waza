package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"fmt"
)

// GetUserByID is the resolver for the getUserById field.
func (r *queryResolver) GetUserByID(ctx context.Context, id string) (*User, error) {
	panic(fmt.Errorf("not implemented: GetUserByID - getUserById"))
}

// GetUserByEmail is the resolver for the getUserByEmail field.
func (r *queryResolver) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	panic(fmt.Errorf("not implemented: GetUserByEmail - getUserByEmail"))
}

// GetUserByPhone is the resolver for the getUserByPhone field.
func (r *queryResolver) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	panic(fmt.Errorf("not implemented: GetUserByPhone - getUserByPhone"))
}

// ListAccounts is the resolver for the listAccounts field.
func (r *queryResolver) ListAccounts(ctx context.Context) ([]*Account, error) {
	panic(fmt.Errorf("not implemented: ListAccounts - listAccounts"))
}

// GetAccountByID is the resolver for the getAccountById field.
func (r *queryResolver) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	panic(fmt.Errorf("not implemented: GetAccountByID - getAccountById"))
}

// GetTransactionHistory is the resolver for the getTransactionHistory field.
func (r *queryResolver) GetTransactionHistory(ctx context.Context, accountID string) ([]*Transaction, error) {
	panic(fmt.Errorf("not implemented: GetTransactionHistory - getTransactionHistory"))
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
