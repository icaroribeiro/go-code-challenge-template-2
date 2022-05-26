package resolver_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHealthCheckResolverUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetHealthCheck() {
	status := "everything is up and running"

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingStatus",
			SetUp: func(t *testing.T) {
				returnArgs = ReturnArgs{
					{nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfItAnErrorOccursWhenGettingTheStatus",
			SetUp: func(t *testing.T) {
				returnArgs = ReturnArgs{
					{customerror.New("failed")},
				}
			},
			WantError:   true,
			ShouldPanic: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckService := new(healthcheckmockservice.Service)
			healthCheckService.On("GetStatus").Return(returnArgs[0]...)
			authService := new(authmockservice.Service)
			userService := new(usermockservice.Service)

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			cl := client.New(srv)

			query := getHealthCheckQuery

			resp := GetHealthCheckResponse{}

			if !tc.WantError {
				cl.MustPost(query, &resp)
				assert.Equal(t, resp.GetHealthCheck.Status, status)
			} else {
				if tc.ShouldPanic {
					mustPostFuncShouldPanic(t, cl.MustPost, query, resp)
				}
			}
		})
	}
}

func mustPostFuncShouldPanic(t *testing.T, f MustPostFunc, query string, resp interface{}) {
	defer func() { recover() }()
	f(query, &resp)
	t.Errorf("It should have panicked.")
}
