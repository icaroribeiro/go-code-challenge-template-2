package auth_test

import (
	"context"
	"testing"

	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	authdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/directive/auth"
	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestFromContext() {
	authDetailsCtxValue := domainentity.Auth{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAnAssociatedValueFromAContext",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainentityfactory.NewAuth(nil)
				ctx = authdirective.NewContext(ctx, authDetailsCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedAuthDetailsCtxValue, ok := authdirective.FromContext(ctx)

			if !tc.WantError {
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.NotEmpty(t, returnedAuthDetailsCtxValue)
				assert.Equal(t, authDetailsCtxValue, returnedAuthDetailsCtxValue)
			}
		})
	}
}
