package resolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	fake "github.com/brianvoe/gofakeit/v5"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/generated"
	dbtrxmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/mockdirective/dbtrx"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestSignIn() {
	driver := "postgres"
	db, _ := NewMockDB(driver)

	dbTrx := &gorm.DB{}

	credentials := securitypkgfactory.NewCredentials(nil)

	opts := []client.Option{}

	tokenString := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningIn",
			SetUp: func(t *testing.T) {
				dbTrx = db

				opts = []client.Option{}
				opts = append(opts, client.Var("input", credentials))
				ctx := context.Background()
				opts = append(opts, AddDBTrxToCtx(ctx, dbTrx))

				tokenString = fake.Word()

				returnArgs = ReturnArgs{
					{tokenString, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetTheDatabaseTransactionFromTheRequestContext",
			SetUp: func(t *testing.T) {
				dbTrx = nil

				opts = []client.Option{}
				opts = append(opts, client.Var("input", credentials))

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
			SetUp: func(t *testing.T) {
				dbTrx = nil

				opts = []client.Option{}
				opts = append(opts, client.Var("input", credentials))
				ctx := context.Background()
				opts = append(opts, AddDBTrxToCtx(ctx, dbTrx))

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenLoggingIn",
			SetUp: func(t *testing.T) {
				dbTrx = db

				opts = []client.Option{}
				opts = append(opts, client.Var("input", credentials))
				ctx := context.Background()
				opts = append(opts, AddDBTrxToCtx(ctx, dbTrx))

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
			authService.On("LogIn", credentials).Return(returnArgs[0]...)
			userService := new(usermockservice.Service)

			res := resolverpkg.New(healthCheckService, authService, userService)

			cfg := generated.Config{Resolvers: res}

			dbTrxDirective := new(dbtrxmockdirective.Directive)
			dbTrxDirective.On("DBTrxMiddleware").Return(MockDirective())

			cfg.Directives.UseDBTrxMiddleware = dbTrxDirective.DBTrxMiddleware()

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					cfg,
				),
			)

			mutation := signInMutation
			resp := SignInMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opts...)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, resp.SignIn.Token)
				assert.Equal(t, tokenString, resp.SignIn.Token)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}
		})
	}
}
