package security_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/security"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func (ts *TestSuite) TestVerifyPasswords() {
	password := ""
	cost := 0
	hashedPassword := ""

	errorType := customerror.NoType

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInComparingThePasswords",
			SetUp: func(t *testing.T) {
				password = fake.Word()
				cost = bcrypt.DefaultCost
				bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				hashedPassword = string(bytes)
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheHashedPasswordIsNotTheHashOfTheGivenPassword",
			SetUp: func(t *testing.T) {
				password = fake.Word()
				otherPassword := fake.Word()
				cost = bcrypt.DefaultCost
				bytes, err := bcrypt.GenerateFromPassword([]byte(otherPassword), cost)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				hashedPassword = string(bytes)

				errorType = customerror.Unauthorized
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheHashedPasswordIsAnEmptyString",
			SetUp: func(t *testing.T) {
				password = fake.Word()
				hashedPassword = ""

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			security := securitypkg.New()

			err := security.VerifyPasswords(hashedPassword, password)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
			}
		})
	}
}
