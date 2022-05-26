package resolver_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/suite"
)

type Case struct {
	Context     string
	SetUp       func(t *testing.T)
	WantError   bool
	ShouldPanic bool
	TearDown    func(t *testing.T)
}

type Cases []Case

type ReturnArgs [][]interface{}

type TestSuite struct {
	suite.Suite
	Cases Cases
}

type MustPostFunc func(query string, response interface{}, options ...client.Option)
