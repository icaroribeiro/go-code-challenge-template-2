package auth

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// IDirective interface is a collection of function signatures that represents the auth's directive contract.
type IDirective interface {
	AuthMiddleware() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error)
	// AuthRenewalMiddleware() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error)
}
