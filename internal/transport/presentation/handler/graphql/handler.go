package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/auth"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
)

type Handler struct {
	Resolver *resolver.Resolver
}

// New is the factory function that encapsulates the implementation related to graphql handler.
func New(healthCheckService healthcheckservice.IService,
	authService authservice.IService, userService userservice.IService) IHandler {
	res := resolver.NewResolver(healthCheckService, authService, userService)

	return &Handler{
		Resolver: res,
	}
}

func (h *Handler) GraphQL() *handler.Server {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: h.Resolver,
			},
		),
	)

	// srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	// 	res := next(ctx)

	// 	dbTrx := &gorm.DB{}
	// 	if dbTrx = dbtrxmiddleware.ForContext(ctx); dbTrx == nil {
	// 		log.Panic("Failed to retrieve database transaction from context")
	// 	}

	// 	if len(res.Errors) > 0 {
	// 		log.Printf("rolling back database transaction due to response with error(s)")
	// 		if err := dbTrx.Rollback().Error; err != nil {
	// 			log.Panicf("database transaction rollback failed: %s", err.Error())
	// 		}
	// 	} else {
	// 		if err := dbTrx.Commit().Error; err != nil {
	// 			log.Panicf("database transaction commit failed: %s", err.Error())
	// 		}
	// 	}

	// 	return res
	// })

	return srv
}
