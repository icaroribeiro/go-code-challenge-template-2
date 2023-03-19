package resolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	fake "github.com/brianvoe/gofakeit/v5"
	authmockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	"github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/generated"
	dbtrxmockdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/dbtrx"
	resolverpkg "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/resolver"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestSignUp() {
	driver := "postgres"
	db, _ := NewMockDB(driver)

	dbTrx := &gorm.DB{}

	credentials := securitypkg.CredentialsFactory(nil)

	opts := []client.Option{}

	tokenString := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningUp",
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
			Context: "ItShouldFailIfAnErrorOccursWhenRegisteringTheCredentials",
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
			authService.On("Register", credentials).Return(returnArgs[0]...)
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

			mutation := signUpMutation
			resp := SignUpMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opts...)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, resp.SignUp.Token)
				assert.Equal(t, tokenString, resp.SignUp.Token)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}
		})
	}
}
