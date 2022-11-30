package auth_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/dgrijalva/jwt-go"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestFetchAuthFromToken() {
	rsaKeys := ts.RSAKeys
	authpkg := authpkg.New(rsaKeys)

	id := uuid.NewV4()
	userID := uuid.NewV4()

	token := &jwt.Token{}

	errorType := customerror.NoType

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInFetchingAuthDataFromAToken",
			SetUp: func(t *testing.T) {
				claims := jwt.MapClaims{
					"auth_id": id.String(),
					"user_id": userID.String(),
				}

				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				assert.NotNil(t, token, "Token is nil")
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheAuthIDFromTokenIsNotAString",
			SetUp: func(t *testing.T) {
				claims := jwt.MapClaims{
					"auth_id": fake.Number(1, 10),
				}

				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				assert.NotNil(t, token, "Token is nil")

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIDFromTokenIsNotAUUIDString",
			SetUp: func(t *testing.T) {
				claims := jwt.MapClaims{
					"auth_id": fake.Word(),
				}

				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				assert.NotNil(t, token, "Token is nil")

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheUserIDFromTokenIsNotAString",
			SetUp: func(t *testing.T) {
				claims := jwt.MapClaims{
					"auth_id": id.String(),
					"user_id": fake.Number(1, 10),
				}

				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				assert.NotNil(t, token, "Token is nil")

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheUserIDFromTokenIsNotAUUIDString",
			SetUp: func(t *testing.T) {
				claims := jwt.MapClaims{
					"auth_id": id.String(),
					"user_id": fake.Word(),
				}

				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				assert.NotNil(t, token, "Token is nil")

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheTokenIsNil",
			SetUp: func(t *testing.T) {
				token = nil

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			auth, err := authpkg.FetchAuthFromToken(token)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.Equal(t, id, auth.ID)
				assert.Equal(t, userID, auth.UserID)
			} else {
				assert.NotNil(t, err, "Predicted error lost")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, auth, "Auth is not empty")
			}
		})
	}
}
