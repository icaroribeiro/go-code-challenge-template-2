package auth_test

import (
	"fmt"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/dgrijalva/jwt-go"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestValidateTokenRenewal() {
	auth := domainentity.Auth{}

	issuedAt := time.Now().Unix()
	expiredAt := time.Now().Unix()

	rsaKeys := ts.RSAKeys
	authpkg := authpkg.New(rsaKeys)

	tokenString := ""
	timeBeforeTokenExpTimeInSec := 60

	err := customerror.New("failed")

	errorType := customerror.NoType

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInValidatingTokenRenewalIfTokenHasExpired",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainentity.Auth{
					ID:     id,
					UserID: userID,
				}

				issuedAt = time.Now().Unix()
				tokenExpTimeInSec := fake.Number(-60, -30)
				duration := time.Second * time.Duration(tokenExpTimeInSec)
				expiredAt = time.Now().Add(duration).Unix()

				tokenString, err = authpkg.CreateToken(auth, tokenExpTimeInSec)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.NotEmpty(t, tokenString, "")
			},
			WantError: false,
		},
		{
			Context: "ItShouldSucceedInValidatingTokenRenewalIfTokenHasNotExpiredButItsExpTimeIsWithinTheTimePriorToTheTimeBeforeTokenExpTime",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainentity.Auth{
					ID:     id,
					UserID: userID,
				}

				issuedAt = time.Now().Unix()
				tokenExpTimeInSec := fake.Number(30, 60)
				duration := time.Second * time.Duration(tokenExpTimeInSec)
				expiredAt = time.Now().Add(duration).Unix()

				tokenString, err = authpkg.CreateToken(auth, tokenExpTimeInSec)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.NotEmpty(t, tokenString, "")
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheTokenIsInvalid",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainentity.Auth{
					ID:     id,
					UserID: userID,
				}

				tokenString = fake.Word()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTokenHasNotExpiredButItsExpTimeIsNotWithinTheTimePriorToTheTimeBeforeTokenExpTime",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainentity.Auth{
					ID:     id,
					UserID: userID,
				}

				tokenExpTimeInSec := fake.Number(300, 600)

				tokenString, err = authpkg.CreateToken(auth, tokenExpTimeInSec)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.NotEmpty(t, tokenString, "")

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			token, err := authpkg.ValidateTokenRenewal(tokenString, timeBeforeTokenExpTimeInSec)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				claims, ok := token.Claims.(jwt.MapClaims)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.Equal(t, auth.ID.String(), claims["auth_id"])
				assert.Equal(t, auth.UserID.String(), claims["user_id"])
				iat, ok := claims["iat"].(float64)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.WithinDuration(t, time.Unix(issuedAt, 0), time.Unix(int64(iat), 0), time.Second)
				exp, ok := claims["exp"].(float64)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.WithinDuration(t, time.Unix(expiredAt, 0), time.Unix(int64(exp), 0), time.Second)
			} else {
				assert.NotNil(t, err, "Predicted error lost")
				assert.Equal(t, errorType, customerror.GetType(err))
			}
		})
	}
}
