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

var refreshTokenMutation = `mutation {
		refreshToken {
			token
		}
	}`

var changePasswordMutation = `mutation($input: Passwords!) {
		changePassword(input: $input) {
			message
		}
	}`

var signOutMutation = `mutation {
		signOut {
			message
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

type RefreshTokenMutationResponse struct {
	RefreshToken struct {
		Token string
	}
}

type ChangePasswordMutationResponse struct {
	ChangePassword struct {
		Message string
	}
}

type SignOutMutationResponse struct {
	SignOut struct {
		Message string
	}
}
