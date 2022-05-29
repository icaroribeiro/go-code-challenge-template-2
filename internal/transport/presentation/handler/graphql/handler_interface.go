package graphql

import "net/http"

// IHandler interface is the graphql's handler contract.
type IHandler interface {
	GraphQL() func(w http.ResponseWriter, r *http.Request)
}
