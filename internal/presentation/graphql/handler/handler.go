package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/auth"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/user"
	authdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/directive/auth"
	dbtrxdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/directive/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/resolver"
)

type Handler struct {
	Cfg generated.Config
}

// New is the factory function that encapsulates the implementation related to graphql handler.
func New(healthCheckService healthcheckservice.IService, authService authservice.IService, userService userservice.IService,
	dbTrxDirective dbtrxdirective.IDirective, authDirective authdirective.IDirective) IHandler {
	res := resolver.New(healthCheckService, authService, userService)

	cfg := generated.Config{Resolvers: res}
	cfg.Directives.UseDBTrxMiddleware = dbTrxDirective.DBTrxMiddleware()
	cfg.Directives.UseAuthMiddleware = authDirective.AuthMiddleware()
	cfg.Directives.UseAuthRenewalMiddleware = authDirective.AuthRenewalMiddleware()

	return &Handler{
		Cfg: cfg,
	}
}

func (h *Handler) GraphQL() func(w http.ResponseWriter, r *http.Request) {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			h.Cfg,
		),
	)

	return srv.ServeHTTP
}
