package resolver_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GetAllUsersQueryResponse struct {
	GetAllUsers []struct {
		ID       string
		Username string
	}
}

var getAllUsersQuery = `query {
		getAllUsers {
			id
			username
		}
	}`

func TestUserResolverSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
