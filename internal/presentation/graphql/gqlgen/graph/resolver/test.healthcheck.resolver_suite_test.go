package resolver_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

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

func TestHealthCheckResolverSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
