package auth_test

import (
	"context"
	"testing"

	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	authdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/directive/auth"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNewContext() {
	authDetailsCtxValue := domainentity.Auth{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingACopyOfAContextWithAnAssociatedValue",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainentity.AuthFactory(nil)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedCtx := authdirective.NewContext(ctx, authDetailsCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedAuthDetailsCtxValue, ok := authdirective.FromContext(returnedCtx)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.Equal(t, authDetailsCtxValue, returnedAuthDetailsCtxValue)
			}
		})
	}
}
