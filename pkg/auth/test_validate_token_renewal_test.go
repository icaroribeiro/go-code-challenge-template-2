package auth_test

import (
	"fmt"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/golang-jwt/jwt"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestValidateTokenRenewal() {
	expiredAt := time.Now().Unix()

	rsaKeys := ts.RSAKeys
	authpkg := authpkg.New(rsaKeys)

	id := uuid.NewV4()
	userID := uuid.NewV4()

	token := &jwt.Token{}

	timeBeforeTokenExpTimeInSec := 60

	errorType := customerror.NoType

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInValidatingTokenRenewalIfTokenHasExpired",
			SetUp: func(t *testing.T) {
				tokenExpTimeInSec := fake.Number(-60, -30)
				duration := time.Second * time.Duration(tokenExpTimeInSec)
				expiredAt = time.Now().Add(duration).Unix()

				claims := jwt.MapClaims{
					"auth_id": id.String(),
					"user_id": userID.String(),
					"exp":     expiredAt,
				}

				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				assert.NotNil(t, token, "Token is nil")
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTokenHasNotExpiredButItsExpTimeIsNotWithinTheTimePriorToTheTimeBeforeTokenExpTime",
			SetUp: func(t *testing.T) {
				tokenExpTimeInSec := fake.Number(300, 600)
				duration := time.Second * time.Duration(tokenExpTimeInSec)
				expiredAt = time.Now().Add(duration).Unix()

				claims := jwt.MapClaims{
					"auth_id": id.String(),
					"user_id": userID.String(),
					"exp":     expiredAt,
				}

				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				assert.NotNil(t, token, "Token is nil")

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedToken, err := authpkg.ValidateTokenRenewal(token, timeBeforeTokenExpTimeInSec)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				claims, ok := returnedToken.Claims.(jwt.MapClaims)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.Equal(t, id.String(), claims["auth_id"])
				assert.Equal(t, userID.String(), claims["user_id"])
				exp, ok := claims["exp"].(int64)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.WithinDuration(t, time.Unix(expiredAt, 0), time.Unix(int64(exp), 0), time.Second)
			} else {
				assert.NotNil(t, err, "Predicted error lost")
				assert.Equal(t, errorType, customerror.GetType(err))
			}
		})
	}
}
