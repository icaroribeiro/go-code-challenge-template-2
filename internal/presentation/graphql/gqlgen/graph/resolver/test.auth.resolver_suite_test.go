package resolver_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

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

func TestAuthResolverSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
