package auth_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestCreateToken() {
	rsaKeys := ts.RSAKeys

	auth := domainentity.Auth{}

	tokenExpTimeInSec := 0

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingAToken",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainentity.Auth{
					ID:     id,
					UserID: userID,
				}

				tokenExpTimeInSec = fake.Number(30, 60)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authpkg := authpkg.New(rsaKeys)

			tokenString, err := authpkg.CreateToken(auth, tokenExpTimeInSec)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.NotEmpty(t, tokenString, "Unexpected empty token")
			}
		})
	}
}
