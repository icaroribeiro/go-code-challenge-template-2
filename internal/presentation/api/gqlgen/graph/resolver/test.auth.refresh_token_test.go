package resolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	fake "github.com/brianvoe/gofakeit/v5"
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	authmockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	"github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/generated"
	authmockdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/auth"
	resolverpkg "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/resolver"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestRefreshToken() {
	dbTrx := &gorm.DB{}
	dbTrx = nil

	tokenString := fake.Word()

	auth := domainentity.AuthFactory(nil)

	opt := func(bd *client.Request) {}

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRefreshingTheToken",
			SetUp: func(t *testing.T) {
				ctx := context.Background()
				opt = AddAuthDetailsToCtx(ctx, auth)

				returnArgs = ReturnArgs{
					{tokenString, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
			SetUp: func(t *testing.T) {
				opt = func(bd *client.Request) {}

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthFromTheRequestContextIsEmpty",
			SetUp: func(t *testing.T) {
				ctx := context.Background()
				opt = AddAuthDetailsToCtx(ctx, domainentity.Auth{})

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenRefreshingTheToken",
			SetUp: func(t *testing.T) {
				ctx := context.Background()
				opt = AddAuthDetailsToCtx(ctx, auth)

				returnArgs = ReturnArgs{
					{"", customerror.New("failed")},
				}
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckService := new(healthcheckmockservice.Service)
			authService := new(authmockservice.Service)
			authService.On("WithDBTrx", dbTrx).Return(authService)
			authService.On("RenewToken", auth).Return(returnArgs[0]...)
			userService := new(usermockservice.Service)

			res := resolverpkg.New(healthCheckService, authService, userService)

			cfg := generated.Config{Resolvers: res}

			authDirective := new(authmockdirective.Directive)
			authDirective.On("AuthRenewalMiddleware").Return(MockDirective())

			cfg.Directives.UseAuthRenewalMiddleware = authDirective.AuthRenewalMiddleware()

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					cfg,
				),
			)

			mutation := refreshTokenMutation
			resp := RefreshTokenMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opt)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, resp.RefreshToken.Token)
				assert.Equal(t, tokenString, resp.RefreshToken.Token)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}
		})
	}
}
