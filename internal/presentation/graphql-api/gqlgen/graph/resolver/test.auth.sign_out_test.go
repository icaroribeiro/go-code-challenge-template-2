package resolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/generated"
	authmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/mockdirective/auth"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestSignOut() {
	dbTrx := &gorm.DB{}
	dbTrx = nil

	auth := domainentityfactory.NewAuth(nil)

	opt := func(bd *client.Request) {}

	message := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningOut",
			SetUp: func(t *testing.T) {
				ctx := context.Background()
				opt = AddAuthDetailsToCtx(ctx, auth)

				message = "you have logged out successfully"

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
			SetUp: func(t *testing.T) {
				opt = func(bd *client.Request) {}

				returnArgs = ReturnArgs{
					{nil},
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
					{nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenSigningOut",
			SetUp: func(t *testing.T) {
				ctx := context.Background()
				opt = AddAuthDetailsToCtx(ctx, auth)

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
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
			authService.On("LogOut", auth.ID.String()).Return(returnArgs[0]...)
			userService := new(usermockservice.Service)

			res := resolverpkg.New(healthCheckService, authService, userService)

			cfg := generated.Config{Resolvers: res}

			authDirective := new(authmockdirective.Directive)
			authDirective.On("AuthMiddleware").Return(MockDirective())

			cfg.Directives.UseAuthMiddleware = authDirective.AuthMiddleware()

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					cfg,
				),
			)

			mutation := signOutMutation
			resp := SignOutMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opt)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, resp.SignOut.Message)
				assert.Equal(t, message, resp.SignOut.Message)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}
		})
	}
}
