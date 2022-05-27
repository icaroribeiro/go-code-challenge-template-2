package graphql_test

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"

	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql"
	authmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/mockdirective/auth"
	dbtrxmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/mockdirective/dbtrx"
	graphqlrouter "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/router/graphql"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestRouterUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestConfigureRoutes() {
	routes := routehttputilpkg.Routes{}

	healthCheckService := new(healthcheckmockservice.Service)
	authService := new(authmockservice.Service)
	userService := new(usermockservice.Service)

	dbTrxDirective := new(dbtrxmockdirective.Directive)
	dbTrxDirective.On("DBTrxMiddleware").Return(MockSchemaDirective())

	authDirective := new(authmockdirective.Directive)
	authDirective.On("AuthMiddleware").Return(MockSchemaDirective())
	authDirective.On("AuthRenewalMiddleware").Return(MockSchemaDirective())

	graphqlHandler := graphqlhandler.New(healthCheckService, authService, userService, dbTrxDirective, authDirective)

	adapters := map[string]adapterhttputilpkg.Adapter{
		"authMiddleware": authmiddlewarepkg.Auth(),
	}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInConfiguringTheRoutes",
			SetUp: func(t *testing.T) {
				routes = routehttputilpkg.Routes{
					routehttputilpkg.Route{
						Name:   "GraphQL",
						Method: http.MethodPost,
						Path:   "/graphql",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(graphqlHandler.GraphQL().ServeHTTP).
							With(adapters["authMiddleware"]),
					},
				}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedRoutes := graphqlrouter.ConfigureRoutes(graphqlHandler, adapters)

			assert.Equal(t, len(routes), len(returnedRoutes))

			for i := range routes {
				assert.Equal(t, routes[i].Name, returnedRoutes[i].Name)
				assert.Equal(t, routes[i].Method, returnedRoutes[i].Method)
				assert.Equal(t, routes[i].Path, returnedRoutes[i].Path)
				handlerFunc1 := runtime.FuncForPC(reflect.ValueOf(routes[i].HandlerFunc).Pointer()).Name()
				handlerFunc2 := runtime.FuncForPC(reflect.ValueOf(returnedRoutes[i].HandlerFunc).Pointer()).Name()
				assert.Equal(t, handlerFunc1, handlerFunc2)
			}
		})
	}
}
