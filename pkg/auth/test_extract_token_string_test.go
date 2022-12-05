package auth_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestExtractTokenString() {
	authHeaderString := ""

	tokenString := ""

	errorType := customerror.NoType

	rsaKeys := ts.RSAKeys
	authpkg := authpkg.New(rsaKeys)

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInExtractingTokenStringIfTokenIsInformed",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				authHeaderString = "Bearer " + tokenString
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailInExtractingTokenStringIfTheAuthHeaderIsAnEmptyString",
			SetUp: func(t *testing.T) {
				authHeaderString = ""

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailInExtractingTokenStringIfTheTokenIsNotInformed",
			SetUp: func(t *testing.T) {
				tokenString = ""

				authHeaderString = "Bearer " + ""

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedTokenString, err := authpkg.ExtractTokenString(authHeaderString)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.Equal(t, tokenString, returnedTokenString)
			} else {
				assert.NotNil(t, err, "Predicted error lost")
				assert.Equal(t, errorType, customerror.GetType(err))
			}
		})
	}
}
