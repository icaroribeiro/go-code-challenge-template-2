package graphql_test

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/suite"
)

type Case struct {
	Context    string
	SetUp      func(t *testing.T)
	StatusCode int
	WantError  bool
	TearDown   func(t *testing.T)
}

type Cases []Case

type ReturnArgs [][]interface{}

type TestSuite struct {
	suite.Suite
	Cases Cases
}

type GetHealthCheckQueryResponse struct {
	GetHealthCheck struct {
		Status string
	}
}

var getHealthCheckQuery = `query {
	getHealthCheck { 
		status
	}
}`

func MockDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		return next(ctx)
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
