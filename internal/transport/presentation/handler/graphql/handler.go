package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/graph/resolver"
)

type Handler struct {
	Resolver *resolver.Resolver
}

// New is the factory function that encapsulates the implementation related to graphql handler.
func New(healthCheckService healthcheckservice.IService) IHandler {
	res := resolver.NewResolver(healthCheckService)

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
