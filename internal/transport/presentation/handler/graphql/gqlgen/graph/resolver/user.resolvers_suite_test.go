package resolver_test

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
