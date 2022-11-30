package server_test

import (
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

type TestSuite struct {
	suite.Suite
	Cases Cases
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
