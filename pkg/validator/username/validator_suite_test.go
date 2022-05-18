package username_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Case struct {
	Context   string
	SetUp     func(t *testing.T)
	Inf       interface{}
	Param     string
	WantError bool
	TearDown  func(t *testing.T)
}

type Cases []Case

type Foo struct{}

type TestSuite struct {
	suite.Suite
	Cases Cases
}
