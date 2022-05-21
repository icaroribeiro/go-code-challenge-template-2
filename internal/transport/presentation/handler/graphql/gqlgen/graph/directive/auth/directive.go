package auth

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/dgrijalva/jwt-go"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/dbtrx"
	"gorm.io/gorm"
)

var authDetailsCtxKey = &contextKey{"auth_details"}

type contextKey struct {
	name string
}

func buildAuth(db *gorm.DB, authN authpkg.IAuth, token *jwt.Token) (domainmodel.Auth, error) {
	auth, err := authN.FetchAuthFromToken(token)
	if err != nil {
		return domainmodel.Auth{}, err
	}

	// Before proceeding is necessary to check if the user who is performing operations is logged
	// based on the authentication details inserted within in the token.
	authAux := domainmodel.Auth{}

	result := db.Find(&authAux, "id=?", auth.ID)
	if result.Error != nil {
		return domainmodel.Auth{}, result.Error
	}

	if authAux.IsEmpty() {
		errorMessage := "you are not logged in, then perform a login to get a token before proceeding"
		return domainmodel.Auth{}, customerror.BadRequest.New(errorMessage)
	}

	if auth.UserID.String() != authAux.UserID.String() {
		errorMessage := "the token's auth_id and user_id are not associated"
		return domainmodel.Auth{}, customerror.BadRequest.New(errorMessage)
	}

	return auth, nil
}

// IsAuthenticated is the function that...
func IsAuthenticated(authN authpkg.IAuth) func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		tokenString := ""

		if tokenString = authmiddlewarepkg.ForContext(ctx); tokenString == "" {
			return nil, customerror.New("failed to get auth_details key from the context of the request")
		}

		token, err := authN.DecodeToken(tokenString)
		if err != nil {
			return nil, err
		}

		dbTrx := &gorm.DB{}

		if dbTrx = dbtrxmiddlewarepkg.ForContext(ctx); dbTrx == nil {
			return nil, err
		}

		auth, err := buildAuth(dbTrx, authN, token)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, authDetailsCtxKey, auth)
		return next(ctx)
	}
}

// CanTokenAlreadyBeRenewed is the function that...
func CanTokenAlreadyBeRenewed(authN authpkg.IAuth, timeBeforeTokenExpTimeInSec int) func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		tokenString := ""

		if tokenString = authmiddlewarepkg.ForContext(ctx); tokenString == "" {
			return nil, customerror.New("failed to get auth_details key from the context of the request")
		}

		token, err := authN.ValidateTokenRenewal(tokenString, timeBeforeTokenExpTimeInSec)
		if err != nil {
			return nil, err
		}

		dbTrx := &gorm.DB{}

		if dbTrx = dbtrxmiddlewarepkg.ForContext(ctx); dbTrx == nil {
			return nil, err
		}

		auth, err := buildAuth(dbTrx, authN, token)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, authDetailsCtxKey, auth)
		return next(ctx)
	}
}

// ForContext is the function that finds the auth_details from the context.
func ForContext(ctx context.Context) domainmodel.Auth {
	raw, _ := ctx.Value(authDetailsCtxKey).(domainmodel.Auth)
	return raw
}
