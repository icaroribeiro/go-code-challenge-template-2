package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
)

type Handler struct {
	Resolver *resolver.Resolver
}

// New is the factory function that encapsulates the implementation related to graphql handler.
func New(healthCheckService healthcheckservice.IService, userService userservice.IService) IHandler {
	res := resolver.NewResolver(healthCheckService, userService)

	return &Handler{
		Resolver: res,
	}
}

func (h *Handler) GraphQL() *handler.Server {
	return handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: h.Resolver,
			},
		),
	)
}
