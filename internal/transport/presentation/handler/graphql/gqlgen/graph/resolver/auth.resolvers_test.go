package resolver_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/DATA-DOG/go-sqlmock"
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/dgrijalva/jwt-go"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	authdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/auth"
	dbtrxdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/model"
	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
	mockauthpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockauth"
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

			mutation := signUpMutation
			resp := SignUpMutationResponse{}

			cl := client.New(srv)
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
			Context: "ItShouldSucceedInSigningIn",
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
			Context: "ItShouldFailIfAnErrorOccursWhenLoggingIn",
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
			authService.On("LogIn", credentials).Return(returnArgs[0]...)
			userService := new(usermockservice.Service)

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			c.Directives.UseDBTrxMiddleware = dbtrxdirectivepkg.DBTrxMiddleware(nil)

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			mutation := signInMutation
			resp := SignInMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opt, addDBTrxCtxValue(ctx, dbTrxCtxValue))

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
	driver := "postgres"
	db, mock := NewMockDB(driver)

	ctx := context.Background()

	tokenString := fake.Word()

	token := &jwt.Token{}

	timeBeforeTokenExpTimeInSec := 0

	auth := domainmodelfactory.NewAuth(nil)

	authDetailsCtxValue := domainmodel.Auth{}

	sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

	args := map[string]interface{}{
		"id":     auth.ID,
		"userID": auth.UserID,
	}

	authDatastore := datastoremodelfactory.NewAuth(args)

	dbTrx := &gorm.DB{}
	dbTrx = nil

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRefreshingTheToken",
			SetUp: func(t *testing.T) {
				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(authDatastore.ID).
					WillReturnRows(rows)

				authDetailsCtxValue = auth

				returnArgs = ReturnArgs{
					{token, nil},
					{auth, nil},
					{tokenString, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenRefreshingTheToken",
			SetUp: func(t *testing.T) {
				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(authDatastore.ID).
					WillReturnRows(rows)

				authDetailsCtxValue = auth

				returnArgs = ReturnArgs{
					{token, nil},
					{auth, nil},
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
			authService.On("RenewToken", authDetailsCtxValue).Return(returnArgs[2]...)
			userService := new(usermockservice.Service)

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			authN := new(mockauthpkg.Auth)
			authN.On("ValidateTokenRenewal", tokenString, timeBeforeTokenExpTimeInSec).Return(returnArgs[0]...)
			authN.On("FetchAuthFromToken", token).Return(returnArgs[1]...)

			c.Directives.UseAuthRenewalMiddleware = authdirectivepkg.AuthRenewalMiddleware(db, authN, timeBeforeTokenExpTimeInSec)

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			mutation := refreshTokenMutation
			resp := RefreshTokenMutationResponse{}

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, addTokenStringCtxValue(ctx, tokenString))

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
