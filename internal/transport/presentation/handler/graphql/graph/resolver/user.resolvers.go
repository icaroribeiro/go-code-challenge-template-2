package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/graph/model"
)

func (r *mutationResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
