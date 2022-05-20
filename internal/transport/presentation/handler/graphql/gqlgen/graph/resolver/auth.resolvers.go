package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/model"

	dbtrxmiddleware "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
	"gorm.io/gorm"
)

func (r *mutationResolver) SignUp(ctx context.Context, credentials security.Credentials) (*model.AuthPayload, error) {
	dbTrx := &gorm.DB{}

	if dbTrx = dbtrxmiddleware.ForContext(ctx); dbTrx == nil {
		return &model.AuthPayload{}, fmt.Errorf("Problem with database")
	}

	token, err := r.AuthService.WithDBTrx(dbTrx).Register(credentials)
	if err != nil {
		return &model.AuthPayload{}, err
	}

	return &model.AuthPayload{Token: token}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
