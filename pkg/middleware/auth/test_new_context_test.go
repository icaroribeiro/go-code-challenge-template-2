package auth_test

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNewContext() {
	tokenCtxValue := &jwt.Token{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingACopyOfAContextWithAnAssociatedValue",
			SetUp: func(t *testing.T) {
				tokenCtxValue = &jwt.Token{
					Valid: true,
				}
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedCtx := authmiddlewarepkg.NewContext(ctx, tokenCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedTokenCtxValue, ok := authmiddlewarepkg.FromContext(returnedCtx)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.Equal(t, tokenCtxValue, returnedTokenCtxValue)
			}
		})
	}
}
