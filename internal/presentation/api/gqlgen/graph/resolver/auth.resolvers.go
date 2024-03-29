package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"

	authdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/directive/auth"
	dbtrxdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/directive/dbtrx"
	"github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/generated"
	presentableentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/presentity"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/security"
)

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, input security.Credentials) (*presentableentity.AuthPayload, error) {
	dbTrx, ok := dbtrxdirective.FromContext(ctx)
	if !ok || dbTrx == nil {
		return &presentableentity.AuthPayload{}, customerror.New("failed to get db_trx_state value from the request context")
	}

	token, err := r.AuthService.WithDBTrx(dbTrx).Register(input)
	if err != nil {
		return &presentableentity.AuthPayload{}, err
	}

	return &presentableentity.AuthPayload{Token: token}, nil
}

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, input security.Credentials) (*presentableentity.AuthPayload, error) {
	dbTrx, ok := dbtrxdirective.FromContext(ctx)
	if !ok || dbTrx == nil {
		return &presentableentity.AuthPayload{}, customerror.New("failed to get db_trx_state value from the request context")
	}

	token, err := r.AuthService.WithDBTrx(dbTrx).LogIn(input)
	if err != nil {
		return &presentableentity.AuthPayload{}, err
	}

	return &presentableentity.AuthPayload{Token: token}, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context) (*presentableentity.AuthPayload, error) {
	auth, ok := authdirective.FromContext(ctx)
	if !ok || auth.IsEmpty() {
		return &presentableentity.AuthPayload{}, customerror.New("failed to get the auth_details value from the request context")
	}

	token, err := r.AuthService.WithDBTrx(nil).RenewToken(auth)
	if err != nil {
		return &presentableentity.AuthPayload{}, err
	}

	return &presentableentity.AuthPayload{Token: token}, nil
}

// ChangePassword is the resolver for the changePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, input security.Passwords) (*presentableentity.InfoPayload, error) {
	auth, ok := authdirective.FromContext(ctx)
	if !ok || auth.IsEmpty() {
		return &presentableentity.InfoPayload{}, customerror.New("failed to get the auth_details value from the request context")
	}

	err := r.AuthService.WithDBTrx(nil).ModifyPassword(auth.UserID.String(), input)
	if err != nil {
		return &presentableentity.InfoPayload{}, err
	}

	return &presentableentity.InfoPayload{
		Message: "the password has been updated successfully",
	}, nil
}

// SignOut is the resolver for the signOut field.
func (r *mutationResolver) SignOut(ctx context.Context) (*presentableentity.InfoPayload, error) {
	auth, ok := authdirective.FromContext(ctx)
	if !ok || auth.IsEmpty() {
		return &presentableentity.InfoPayload{}, customerror.New("failed to get the auth_details value from the request context")
	}

	err := r.AuthService.WithDBTrx(nil).LogOut(auth.ID.String())
	if err != nil {
		return &presentableentity.InfoPayload{}, err
	}

	return &presentableentity.InfoPayload{
		Message: "you have logged out successfully",
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
