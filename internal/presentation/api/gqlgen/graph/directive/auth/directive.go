package auth

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang-jwt/jwt"
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	authpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/auth"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/middleware/auth"
	"gorm.io/gorm"
)

type Directive struct {
	DB                          *gorm.DB
	AuthN                       authpkg.IAuth
	TimeBeforeTokenExpTimeInSec int
}

// New is the factory function that encapsulate the implementation related to auth directive.
func New(db *gorm.DB, authN authpkg.IAuth, timeBeforeTokenExpTimeInSec int) IDirective {
	return &Directive{
		DB:                          db,
		AuthN:                       authN,
		TimeBeforeTokenExpTimeInSec: timeBeforeTokenExpTimeInSec,
	}
}

var authDetailsCtxKey = &contextKey{"auth_details"}

type contextKey struct {
	name string
}

func buildAuth(db *gorm.DB, authN authpkg.IAuth, token *jwt.Token) (domainentity.Auth, error) {
	auth, err := authN.FetchAuthFromToken(token)
	if err != nil {
		return domainentity.Auth{}, err
	}

	// Before proceeding is necessary to check if the user who is performing operations is logged
	// based on the authentication details inserted within in the token.
	authAux := domainentity.Auth{}

	result := db.Find(&authAux, "id=?", auth.ID)
	if result.Error != nil {
		return domainentity.Auth{}, result.Error
	}

	if authAux.IsEmpty() {
		errorMessage := "you are not logged in, then perform a login to get a token before proceeding"
		return domainentity.Auth{}, customerror.BadRequest.New(errorMessage)
	}

	if auth.UserID.String() != authAux.UserID.String() {
		errorMessage := "the token's auth_id and user_id are not associated"
		return domainentity.Auth{}, customerror.BadRequest.New(errorMessage)
	}

	return auth, nil
}

// AuthMiddleware is the function that acts as a HTTP middleware to evaluate the authentication of API based on a JWT token.
func (d *Directive) AuthMiddleware() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		token, ok := authmiddlewarepkg.FromContext(ctx)
		if !ok {
			return nil, customerror.New("failed to get the token value from the request context")
		}

		auth, err := buildAuth(d.DB, d.AuthN, token)
		if err != nil {
			return nil, err
		}

		ctx = NewContext(ctx, auth)

		return next(ctx)
	}
}

// AuthRenewalMiddleware is the function that acts as a HTTP middleware to evaluate the authentication renewal of API based on a JWT token.
func (d *Directive) AuthRenewalMiddleware() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		token, ok := authmiddlewarepkg.FromContext(ctx)
		if !ok {
			return nil, customerror.New("failed to get the token value from the request context")
		}

		token, err := d.AuthN.ValidateTokenRenewal(token, d.TimeBeforeTokenExpTimeInSec)
		if err != nil {
			return nil, err
		}

		auth, err := buildAuth(d.DB, d.AuthN, token)
		if err != nil {
			return nil, err
		}

		ctx = NewContext(ctx, auth)

		return next(ctx)
	}
}

// NewContext is the function that returns a new Context that carries auth_details value.
func NewContext(ctx context.Context, auth domainentity.Auth) context.Context {
	return context.WithValue(ctx, authDetailsCtxKey, auth)
}

// FromContext is the function that returns the auth_details value stored in context, if any.
func FromContext(ctx context.Context) (domainentity.Auth, bool) {
	raw, ok := ctx.Value(authDetailsCtxKey).(domainentity.Auth)
	return raw, ok
}
