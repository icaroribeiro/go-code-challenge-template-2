package auth_test

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestFromContext() {
	tokenCtxValue := &jwt.Token{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAnAssociatedValueFromAContext",
			SetUp: func(t *testing.T) {
				tokenCtxValue = &jwt.Token{
					Valid: true,
				}
				ctx = authmiddlewarepkg.NewContext(ctx, tokenCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedTokenCtxValue, ok := authmiddlewarepkg.FromContext(ctx)

			if !tc.WantError {
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.NotEmpty(t, returnedTokenCtxValue)
				assert.Equal(t, tokenCtxValue, returnedTokenCtxValue)
			}
		})
	}
}
