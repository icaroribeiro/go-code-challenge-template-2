package resolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	fake "github.com/brianvoe/gofakeit/v5"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
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

	dbTrx := &gorm.DB{}

	tokenString := ""

	ctx := context.Background()

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningUp",
			SetUp: func(t *testing.T) {
				dbTrx = db

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
				dbTrx = nil

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

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			c.Directives.UseDBTrxMiddleware = MockSchemaDirective()

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			mutation := signUpMutation
			resp := SignUpMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opt, addDBTrxToCtx(ctx, dbTrx))

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

	dbTrx := &gorm.DB{}

	tokenString := ""

	ctx := context.Background()

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningIn",
			SetUp: func(t *testing.T) {
				dbTrx = db

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
				dbTrx = nil

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

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			c.Directives.UseDBTrxMiddleware = MockSchemaDirective()

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			mutation := signInMutation
			resp := SignInMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opt, addDBTrxToCtx(ctx, dbTrx))

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

func (ts *TestSuite) TestRefreshToken() {
	ctx := context.Background()

	tokenString := fake.Word()

	auth := domainmodelfactory.NewAuth(nil)

	authDetails := domainmodel.Auth{}

	dbTrx := &gorm.DB{}
	dbTrx = nil

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRefreshingTheToken",
			SetUp: func(t *testing.T) {
				authDetails = auth

				returnArgs = ReturnArgs{
					{tokenString, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
			SetUp: func(t *testing.T) {
				authDetails = domainmodel.Auth{}

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenRefreshingTheToken",
			SetUp: func(t *testing.T) {
				authDetails = auth

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
			authService.On("RenewToken", authDetails).Return(returnArgs[0]...)
			userService := new(usermockservice.Service)

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			// authN := new(mockauthpkg.Auth)
			// authN.On("ValidateTokenRenewal", tokenString, timeBeforeTokenExpTimeInSec).Return(returnArgs[0]...)
			// authN.On("FetchAuthFromToken", token).Return(returnArgs[1]...)

			c.Directives.UseAuthRenewalMiddleware = MockSchemaDirective()

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			mutation := refreshTokenMutation
			resp := RefreshTokenMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, addAuthDetailsToCtx(ctx, authDetails))

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
