package dbtrx

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// IDirective interface is a collection of function signatures that represents the dbtrx's directive contract.
type IDirective interface {
	DBTrxMiddleware() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error)
}
