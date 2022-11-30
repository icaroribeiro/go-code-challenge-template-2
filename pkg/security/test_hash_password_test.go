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

func (ts *TestSuite) TestHashPassword() {
	password := ""
	cost := 0

	errorType := customerror.NoType

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInHashingThePassword",
			SetUp: func(t *testing.T) {
				password = fake.Word()
				cost = bcrypt.DefaultCost
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfThePasswordIsAnEmptyString",
			SetUp: func(t *testing.T) {
				password = ""
				cost = bcrypt.DefaultCost

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailInHashingThePasswordIfTheCostIsBiggerThan31",
			SetUp: func(t *testing.T) {
				password = fake.Word()
				cost = fake.Number(32, 100)

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			security := securitypkg.New()

			bytes, err := security.HashPassword(password, cost)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				hashedPassword := string(bytes)
				err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, bytes)
			}
		})
	}
}
