package resolver_test

var getAllUsersQuery = `query {
		getAllUsers {
			id
			username
		}
	}`

type GetAllUsersQueryResponse struct {
	GetAllUsers []struct {
		ID       string
		Username string
	}
}
