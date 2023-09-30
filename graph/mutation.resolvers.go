package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"fmt"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, payload CreateUserInput) (*Result, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// TransferFunds is the resolver for the transferFunds field.
func (r *mutationResolver) TransferFunds(ctx context.Context, payload *TransferFundsInput) (*Result, error) {
	panic(fmt.Errorf("not implemented: TransferFunds - transferFunds"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }