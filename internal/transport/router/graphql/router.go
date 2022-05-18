package graphql

import (
	"net/http"

	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/route"
)

// ConfigureRoutes is the function that arranges the graphql's routes.
func ConfigureRoutes(graphqlHandler graphqlhandler.IHandler) routehttputilpkg.Routes {
	return routehttputilpkg.Routes{
		routehttputilpkg.Route{
			Name:        "GraphQL",
			Method:      http.MethodPost,
			Path:        "/graphql",
			HandlerFunc: graphqlHandler.GraphQL().ServeHTTP,
		},
	}
}
