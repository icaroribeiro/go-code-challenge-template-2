package graphql

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
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
	res := resolver.NewResolver(healthCheckService, authService, userService)

	return &Handler{
		Resolver:                    res,
		DB:                          db,
		AuthN:                       authN,
		TimeBeforeTokenExpTimeInSec: timeBeforeTokenExpTimeInSec,
	}
}

func (h *Handler) GraphQL() *handler.Server {
	c := generated.Config{Resolvers: h.Resolver}

	c.Directives.UseDBTrx = dbtrxdirectivepkg.UseDBTrx(h.DB)
	c.Directives.IsAuthenticated = authdirectivepkg.IsAuthenticated(h.DB, h.AuthN)
	c.Directives.CanTokenAlreadyBeRenewed = authdirectivepkg.CanTokenAlreadyBeRenewed(h.DB, h.AuthN, h.TimeBeforeTokenExpTimeInSec)

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			c,
		),
	)

	//srv.AroundOperations()

	srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		log.Println(">>> 1")
		dbTrx, ok := dbtrxdirectivepkg.FromContext(ctx)
		if !ok || dbTrx == nil {
			log.Println("Não tem banco!")
		}

		res := next(ctx)

		log.Println(">>> 2")
		dbTrx, ok = dbtrxdirectivepkg.FromContext(ctx)
		if !ok || dbTrx == nil {
			log.Println("Não tem banco após o res!")
		}

		// if !dbTrxState.WasChanged {
		// 	log.Println("Não alterou o estado!")
		// 	if err := dbTrxState.DBTrx.Rollback().Error; err != nil {
		// 		log.Panicf("failed to rollback database transaction: %s", err.Error())
		// 	}
		// 	return res
		// }

		// log.Println("Alterou o estado!")
		// if len(res.Errors) > 0 {
		// 	log.Printf("rolling back database transaction due to error(s)")
		// 	if err := dbTrx.Rollback().Error; err != nil {
		// 		log.Panicf("failed to rollback database transaction: %s", err.Error())
		// 	}
		// } else {
		// 	if err := dbTrx.Commit().Error; err != nil {
		// 		log.Panicf("failed to commit database transaction: %s", err.Error())
		// 	}
		// }

		return res
	})

	return srv
}
