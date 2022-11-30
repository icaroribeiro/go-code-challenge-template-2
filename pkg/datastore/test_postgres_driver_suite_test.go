package datastore_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestPostgresDriverSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
