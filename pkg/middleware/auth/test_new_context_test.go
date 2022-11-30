package auth_test

import (
	"context"
	"testing"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/core/domain/entity"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNewContext() {
	authDetailsCtxValue := domainmodel.Auth{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingACopyOfAContextWithAnAssociatedValue",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainmodelfactory.NewAuth(nil)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedCtx := authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedAuthDetailsCtxValue, ok := authmiddlewarepkg.FromContext(returnedCtx)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.Equal(t, authDetailsCtxValue, returnedAuthDetailsCtxValue)
			}
		})
	}
}
