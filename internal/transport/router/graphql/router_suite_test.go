package graphql_test

import (
	"net/http"
	"testing"

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

func GraphqlHandler(w http.ResponseWriter, r *http.Request) {}
