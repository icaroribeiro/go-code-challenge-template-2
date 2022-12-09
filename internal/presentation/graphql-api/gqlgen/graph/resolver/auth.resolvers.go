package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	authdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/directive/auth"
	dbtrxdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/directive/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/entity"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
)

func (r *mutationResolver) SignUp(ctx context.Context, input security.Credentials) (*entity.AuthPayload, error) {
	dbTrx, ok := dbtrxdirective.FromContext(ctx)
	if !ok || dbTrx == nil {
		return &entity.AuthPayload{}, customerror.New("failed to get db_trx_state value from the request context")
	}

	token, err := r.AuthService.WithDBTrx(dbTrx).Register(input)
	if err != nil {
		return &entity.AuthPayload{}, err
	}

	return &entity.AuthPayload{Token: token}, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, input security.Credentials) (*entity.AuthPayload, error) {
	dbTrx, ok := dbtrxdirective.FromContext(ctx)
	if !ok || dbTrx == nil {
		return &entity.AuthPayload{}, customerror.New("failed to get db_trx_state value from the request context")
	}

	token, err := r.AuthService.WithDBTrx(dbTrx).LogIn(input)
	if err != nil {
		return &entity.AuthPayload{}, err
	}

	return &entity.AuthPayload{Token: token}, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context) (*entity.AuthPayload, error) {
	auth, ok := authdirective.FromContext(ctx)
	if !ok || auth.IsEmpty() {
		return &entity.AuthPayload{}, customerror.New("failed to get the auth_details value from the request context")
	}

	token, err := r.AuthService.WithDBTrx(nil).RenewToken(auth)
	if err != nil {
		return &entity.AuthPayload{}, err
	}

	return &entity.AuthPayload{Token: token}, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input security.Passwords) (*entity.InfoPayload, error) {
	auth, ok := authdirective.FromContext(ctx)
	if !ok || auth.IsEmpty() {
		return &entity.InfoPayload{}, customerror.New("failed to get the auth_details value from the request context")
	}

	err := r.AuthService.WithDBTrx(nil).ModifyPassword(auth.UserID.String(), input)
	if err != nil {
		return &entity.InfoPayload{}, err
	}

	return &entity.InfoPayload{
		Message: "the password has been updated successfully",
	}, nil
}

func (r *mutationResolver) SignOut(ctx context.Context) (*entity.InfoPayload, error) {
	auth, ok := authdirective.FromContext(ctx)
	if !ok || auth.IsEmpty() {
		return &entity.InfoPayload{}, customerror.New("failed to get the auth_details value from the request context")
	}

	err := r.AuthService.WithDBTrx(nil).LogOut(auth.ID.String())
	if err != nil {
		return &entity.InfoPayload{}, err
	}

	return &entity.InfoPayload{
		Message: "you have logged out successfully",
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
