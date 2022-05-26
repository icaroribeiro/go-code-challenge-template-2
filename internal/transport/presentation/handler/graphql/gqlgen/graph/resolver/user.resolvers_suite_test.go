package resolver_test

var getAllUsersQuery = `query {
		getAllUsers {
			id
			username
		}
	}`

type GetAllUsersResponse struct {
	GetAllUsers []struct {
		ID       string
		Username string
	}
}
