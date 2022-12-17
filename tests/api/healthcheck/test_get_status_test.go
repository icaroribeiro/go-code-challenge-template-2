package healthcheck_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/client"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/healthcheck"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	authmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/auth"
	dbtrxmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/dbtrx"
	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/handler"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestGetStatus() {
	db := &gorm.DB{}

	status := ""

	var connPool gorm.ConnPool

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingTheStatus",
			SetUp: func(t *testing.T) {
				db = ts.DB

				status = "everything is up and running"
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseConnectionPoolIsInvalid",
			SetUp: func(t *testing.T) {
				connPool = ts.DB.ConnPool
				ts.DB.ConnPool = nil
				db = ts.DB
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				ts.DB.ConnPool = connPool
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckService := healthcheckservice.New(db)
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
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}

			tc.TearDown(t)
		})
	}
}
