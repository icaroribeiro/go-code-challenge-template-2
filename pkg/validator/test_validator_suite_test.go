package validator_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Case struct {
	Context   string
	SetUp     func(t *testing.T)
	Inf       interface{}
	Tags      string
	WantError bool
	TearDown  func(t *testing.T)
}

type Cases []Case

type Foo struct {
	Field1 string `validate:"nonzero"`
}

type TestSuite struct {
	suite.Suite
	Cases Cases
}

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
