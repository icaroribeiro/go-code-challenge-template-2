package graphql_test

// import (
// 	"net/http"
// 	"reflect"
// 	"runtime"
// 	"testing"

// 	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
// 	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
// 	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql"
// 	graphqlrouter "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/router/graphql"
// 	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/route"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// func TestRouterUnit(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestConfigureRoutes() {
// 	routes := routehttputilpkg.Routes{}

// 	healthCheckService := new(healthcheckmockservice.Service)
// 	userService := new(usermockservice.Service)

// 	graphqlHandler := graphqlhandler.New(healthCheckService, userService)

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInConfiguringTheRoutes",
// 			SetUp: func(t *testing.T) {
// 				routes = routehttputilpkg.Routes{
// 					routehttputilpkg.Route{
// 						Name:        "GraphQL",
// 						Method:      http.MethodPost,
// 						Path:        "/graphql",
// 						HandlerFunc: graphqlHandler.GraphQL().ServeHTTP,
// 					},
// 				}
// 			},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			returnedRoutes := graphqlrouter.ConfigureRoutes(graphqlHandler)

// 			assert.Equal(t, len(routes), len(returnedRoutes))

// 			for i := range routes {
// 				assert.Equal(t, routes[i].Name, returnedRoutes[i].Name)
// 				assert.Equal(t, routes[i].Method, returnedRoutes[i].Method)
// 				assert.Equal(t, routes[i].Path, returnedRoutes[i].Path)
// 				handlerFunc1 := runtime.FuncForPC(reflect.ValueOf(routes[i].HandlerFunc).Pointer()).Name()
// 				handlerFunc2 := runtime.FuncForPC(reflect.ValueOf(returnedRoutes[i].HandlerFunc).Pointer()).Name()
// 				assert.Equal(t, handlerFunc1, handlerFunc2)
// 			}
// 		})
// 	}
// }
