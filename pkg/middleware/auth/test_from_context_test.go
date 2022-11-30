package auth_test

import (
	"context"
	"testing"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/core/domain/entity"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestFromContext() {
	authDetailsCtxValue := domainmodel.Auth{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAnAssociatedValueFromAContext",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainmodelfactory.NewAuth(nil)
				ctx = authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedAuthDetailsCtxValue, ok := authmiddlewarepkg.FromContext(ctx)

			if !tc.WantError {
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.NotEmpty(t, returnedAuthDetailsCtxValue)
				assert.Equal(t, authDetailsCtxValue, returnedAuthDetailsCtxValue)
			}
		})
	}
}
