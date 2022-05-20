package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/model"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
	"gorm.io/gorm"
)

func (r *mutationResolver) SignUp(ctx context.Context, input security.Credentials) (*model.AuthPayload, error) {
	dbTrx := &gorm.DB{}

	if dbTrx = dbtrxmiddlewarepkg.ForContext(ctx); dbTrx == nil {
		return &model.AuthPayload{}, customerror.New("failed to get db_trx key from the context of the request")
	}

	token, err := r.AuthService.WithDBTrx(dbTrx).Register(input)
	if err != nil {
		return &model.AuthPayload{}, err
	}

	return &model.AuthPayload{Token: token}, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, input security.Credentials) (*model.AuthPayload, error) {
	dbTrx := &gorm.DB{}

	if dbTrx = dbtrxmiddlewarepkg.ForContext(ctx); dbTrx == nil {
		return &model.AuthPayload{}, customerror.New("failed to get db_trx key from the context of the request")
	}

	token, err := r.AuthService.WithDBTrx(dbTrx).LogIn(input)
	if err != nil {
		return &model.AuthPayload{}, err
	}

	return &model.AuthPayload{Token: token}, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context) (*model.AuthPayload, error) {
	auth := domainmodel.Auth{}

	if auth = authdirectivepkg.ForContext(ctx); auth.IsEmpty() {
		return &model.AuthPayload{}, customerror.New("failed to get auth_details key from the context of the request")
	}

	token, err := r.AuthService.WithDBTrx(nil).RenewToken(auth)
	if err != nil {
		return &model.AuthPayload{}, err
	}

	return &model.AuthPayload{Token: token}, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input security.Passwords) (*model.InfoPayload, error) {
	auth := domainmodel.Auth{}

	if auth = authdirectivepkg.ForContext(ctx); auth.IsEmpty() {
		return &model.InfoPayload{}, customerror.New("failed to get auth_details key from the context of the request")
	}

	err := r.AuthService.WithDBTrx(nil).ModifyPassword(auth.UserID.String(), input)
	if err != nil {
		return &model.InfoPayload{}, err
	}

	return &model.InfoPayload{
		Message: "the password has been updated successfully",
	}, nil
}

func (r *mutationResolver) SignOut(ctx context.Context) (*model.InfoPayload, error) {
	auth := domainmodel.Auth{}

	if auth = authdirectivepkg.ForContext(ctx); auth.IsEmpty() {
		return &model.InfoPayload{}, customerror.New("failed to get auth_details key from the context of the request")
	}

	err := r.AuthService.WithDBTrx(nil).LogOut(auth.ID.String())
	if err != nil {
		return &model.InfoPayload{}, err
	}

	return &model.InfoPayload{
		Message: "you have logged out successfully",
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
