package router_test

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"

	authmockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	authmockdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/auth"
	dbtrxmockdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/dbtrx"
	graphqlhandler "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/handler"
	graphqlrouter "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/router"
	adapterhttputilpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/middleware/auth"
	mockauthpkg "github.com/icaroribeiro/go-code-challenge-template-2/tests/mocks/pkg/mockauth"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestConfigureRoutes() {
	authN := new(mockauthpkg.Auth)

	routes := routehttputilpkg.Routes{}

	healthCheckService := new(healthcheckmockservice.Service)
	authService := new(authmockservice.Service)
	userService := new(usermockservice.Service)

	dbTrxDirective := new(dbtrxmockdirective.Directive)
	dbTrxDirective.On("DBTrxMiddleware").Return(MockDirective())

	authDirective := new(authmockdirective.Directive)
	authDirective.On("AuthMiddleware").Return(MockDirective())
	authDirective.On("AuthRenewalMiddleware").Return(MockDirective())

	graphqlHandler := graphqlhandler.New(healthCheckService, authService, userService, dbTrxDirective, authDirective)

	adapters := map[string]adapterhttputilpkg.Adapter{
		"authMiddleware": authmiddlewarepkg.Auth(authN),
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
						HandlerFunc: adapterhttputilpkg.AdaptFunc(graphqlHandler.GraphQL()).
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
