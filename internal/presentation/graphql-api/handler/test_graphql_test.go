package graphql_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/client"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	authmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/mockdirective/auth"
	dbtrxmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/mockdirective/dbtrx"
	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/handler"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestGraphQL() {
	status := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInStartingGraphQLServer",
			SetUp: func(t *testing.T) {
				status = "everything is up and running"

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckService := new(healthcheckmockservice.Service)
			healthCheckService.On("GetStatus").Return(returnArgs[0]...)
			authService := new(authmockservice.Service)
			userService := new(usermockservice.Service)

			dbTrxDirective := new(dbtrxmockdirective.Directive)
			dbTrxDirective.On("DBTrxMiddleware").Return(MockDirective())

			authDirective := new(authmockdirective.Directive)
			authDirective.On("AuthMiddleware").Return(MockDirective())
			authDirective.On("AuthRenewalMiddleware").Return(MockDirective())

			graphqlHandler := graphqlhandler.New(healthCheckService, authService, userService, dbTrxDirective, authDirective)

			srv := http.HandlerFunc(graphqlHandler.GraphQL())

			query := getHealthCheckQuery
			resp := GetHealthCheckQueryResponse{}

			cl := client.New(srv)
			err := cl.Post(query, &resp)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, resp.GetHealthCheck.Status, status)
			}
		})
	}
}
