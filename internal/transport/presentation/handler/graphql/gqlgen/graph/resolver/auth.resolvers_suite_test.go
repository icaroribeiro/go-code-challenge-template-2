package resolver_test

var signUpMutation = `mutation($input: Credentials!) {
		signUp(input: $input) {
			token
		}
	}`

var signInMutation = `mutation($input: Credentials!) {
		signIn(input: $input) {
			token
		}
	}`

type SignUpMutationResponse struct {
	SignUp struct {
		Token string
	}
}

type SignInMutationResponse struct {
	SignIn struct {
		Token string
	}
}
