package resolver_test

import (
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/generated"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestGetHealthCheck() {
	status := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingStatus",
			SetUp: func(t *testing.T) {
				status = "everything is up and running"

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfItAnErrorOccursWhenGettingTheStatus",
			SetUp: func(t *testing.T) {
				status = ""

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

			res := resolverpkg.New(healthCheckService, authService, userService)

			cfg := generated.Config{Resolvers: res}

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					cfg,
				),
			)

			query := getHealthCheckQuery
			resp := GetHealthCheckQueryResponse{}

			cl := client.New(srv)
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
