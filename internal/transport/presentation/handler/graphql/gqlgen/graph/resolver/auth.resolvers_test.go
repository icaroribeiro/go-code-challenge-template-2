package resolver_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	fake "github.com/brianvoe/gofakeit/v5"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	dbtrxdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestAuthResolversUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestSignUp() {
	credentials := securitypkgfactory.NewCredentials(nil)

	opt := client.Var("input", credentials)

	driver := "postgres"
	db, _ := NewMockDB(driver)

	dbTrxCtxValue := &gorm.DB{}

	tokenString := ""

	ctx := context.Background()

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningUp",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db

				tokenString = fake.Word()

				returnArgs = ReturnArgs{
					{tokenString, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = nil

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenRegisteringTheCredentials",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db

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
			authService.On("WithDBTrx", dbTrxCtxValue).Return(authService)
			authService.On("Register", credentials).Return(returnArgs[0]...)
			userService := new(usermockservice.Service)

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			c.Directives.UseDBTrxMiddleware = dbtrxdirectivepkg.DBTrxMiddleware(nil)

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			cl := client.New(srv)

			mutation := signUpMutation

			resp := SignUpMutationResponse{}

			err := cl.Post(mutation, &resp, opt, addDBTrxCtxValue(ctx, dbTrxCtxValue))

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

func (ts *TestSuite) TestSignIn() {
	credentials := securitypkgfactory.NewCredentials(nil)

	opt := client.Var("input", credentials)

	driver := "postgres"
	db, _ := NewMockDB(driver)

	dbTrxCtxValue := &gorm.DB{}

	tokenString := ""

	ctx := context.Background()

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInLoggingIn",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				dbTrxCtxValue = db

				tokenString = fake.Word()

				returnArgs = ReturnArgs{
					{tokenString, nil},
				}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				dbTrxCtxValue = nil

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
					"username":"%s",
					"password":"%s"
				`,
					credentials.Username, credentials.Password)

				dbTrxCtxValue = db

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			StatusCode: http.StatusBadRequest,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenLoggingIn",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				dbTrxCtxValue = db

				returnArgs = ReturnArgs{
					{"", customerror.New("failed")},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckService := new(healthcheckmockservice.Service)
			authService := new(authmockservice.Service)
			authService.On("WithDBTrx", dbTrxCtxValue).Return(authService)
			authService.On("Register", credentials).Return(returnArgs[0]...)
			userService := new(usermockservice.Service)

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			c.Directives.UseDBTrxMiddleware = dbtrxdirectivepkg.DBTrxMiddleware(nil)

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			cl := client.New(srv)

			mutation := signUpMutation

			resp := SignUpMutationResponse{}

			err := cl.Post(mutation, &resp, opt, addDBTrxCtxValue(ctx, dbTrxCtxValue))

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
