package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/auth"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/user"
	authdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/auth"
	dbtrxdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	"gorm.io/gorm"
)

type Handler struct {
	Resolver                    *resolver.Resolver
	DB                          *gorm.DB
	AuthN                       authpkg.IAuth
	TimeBeforeTokenExpTimeInSec int
}

// New is the factory function that encapsulates the implementation related to graphql handler.
func New(healthCheckService healthcheckservice.IService, authService authservice.IService, userService userservice.IService,
	db *gorm.DB, authN authpkg.IAuth, timeBeforeTokenExpTimeInSec int) IHandler {
	res := resolver.New(healthCheckService, authService, userService)

	return &Handler{
		Resolver:                    res,
		DB:                          db,
		AuthN:                       authN,
		TimeBeforeTokenExpTimeInSec: timeBeforeTokenExpTimeInSec,
	}
}

func (h *Handler) GraphQL() *handler.Server {
	c := generated.Config{Resolvers: h.Resolver}

	c.Directives.UseDBTrxMiddleware = dbtrxdirectivepkg.DBTrxMiddleware(h.DB)
	c.Directives.UseAuthMiddleware = authdirectivepkg.AuthMiddleware(h.DB, h.AuthN)
	c.Directives.UseAuthRenewalMiddleware = authdirectivepkg.AuthRenewalMiddleware(h.DB, h.AuthN, h.TimeBeforeTokenExpTimeInSec)

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			c,
		),
	)

	return srv
}
