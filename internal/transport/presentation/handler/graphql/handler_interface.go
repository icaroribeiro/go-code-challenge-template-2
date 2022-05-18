package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
)

// IHandler interface is the graphql's handler contract.
type IHandler interface {
	GraphQL() *handler.Server
}
