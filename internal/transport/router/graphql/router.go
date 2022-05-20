package graphql

import (
	"net/http"

	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/route"
)

// ConfigureRoutes is the function that arranges the graphql's routes.
func ConfigureRoutes(graphqlHandler graphqlhandler.IHandler, adapters map[string]adapterhttputilpkg.Adapter) routehttputilpkg.Routes {
	return routehttputilpkg.Routes{
		routehttputilpkg.Route{
			Name:   "GraphQL",
			Method: http.MethodPost,
			Path:   "/graphql",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(graphqlHandler.GraphQL().ServeHTTP).With(
				adapters["authMiddleware"],
				adapters["dbTrxMiddleware"],
			),
		},
	}
}
