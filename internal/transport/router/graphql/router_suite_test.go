package graphql_test

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/suite"
)

type Case struct {
	Context   string
	SetUp     func(t *testing.T)
	WantError bool
	TearDown  func(t *testing.T)
}

type Cases []Case

type ReturnArgs [][]interface{}

type TestSuite struct {
	suite.Suite
	Cases Cases
}

func MockSchemaDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		return next(ctx)
	}
}
