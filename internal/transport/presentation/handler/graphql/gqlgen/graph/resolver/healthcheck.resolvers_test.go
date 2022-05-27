package resolver_test

import (
	"fmt"
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

func TestHealthCheckResolversUnit(t *testing.T) {
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
			WantError: true,
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

			resp := GetHealthCheckQueryResponse{}

			err := cl.Post(query, &resp)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, resp.GetHealthCheck.Status, status)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}
		})
	}
}
