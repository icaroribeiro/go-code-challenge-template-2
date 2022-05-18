package swagger_test

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"

	swaggerrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/router/swagger"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	loggingmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	httpswaggerpkg "github.com/swaggo/http-swagger"
)

func TestRouterUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestConfigureRoutes() {
	routes := routehttputilpkg.Routes{}

	swaggerHandler := httpswaggerpkg.WrapHandler

	adapters := map[string]adapterhttputilpkg.Adapter{
		"loggingMiddleware": loggingmiddlewarepkg.Logging(),
	}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInConfiguringTheRoutes",
			SetUp: func(t *testing.T) {
				routes = routehttputilpkg.Routes{
					routehttputilpkg.Route{
						Name:       "Swagger",
						Method:     http.MethodGet,
						PathPrefix: "/swagger",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(swaggerHandler).
							With(adapters["loggingMiddleware"]),
					},
				}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedRoutes := swaggerrouter.ConfigureRoutes(swaggerHandler, adapters)

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
