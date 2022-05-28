package auth_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/dgrijalva/jwt-go"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/model"
	mockauthpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockauth"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDirectiveUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestNewContext() {
	authDetailsCtxValue := domainmodel.Auth{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingACopyOfAContextWithAnAssociatedValue",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainmodelfactory.NewAuth(nil)
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

func (ts *TestSuite) TestFromContext() {
	authDetailsCtxValue := domainmodel.Auth{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAnAssociatedValueFromAContext",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainmodelfactory.NewAuth(nil)
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

func (ts *TestSuite) TestAuthMiddleware() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	ctx := context.Background()

	tokenString := ""

	token := &jwt.Token{}

	next := func(ctx context.Context) (interface{}, error) { return nil, nil }

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithAuthenticationMiddleware",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     id,
					"userID": userID,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				authDatastore := datastoremodelfactory.NewAuth(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheAuthenticationTokenIsNotSetInTheContext",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				returnArgs = ReturnArgs{
					{nil, nil},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheTokenIsNotDecoded",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				returnArgs = ReturnArgs{
					{nil, customerror.New("failed")},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFetchedFromTheToken",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodel.Auth{}, customerror.New("failed")},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenTryingToFindTheAuthInTheDatabase",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnError(customerror.New("failed"))
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFoundInTheDatabase",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(&sqlmock.Rows{})
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheUserIDFromTokenDoesNotMatchWithTheUserIDFromAuthRecordFromTheDatabase",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				authDatastore := datastoremodelfactory.NewAuth(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authN := new(mockauthpkg.Auth)
			authN.On("DecodeToken", tokenString).Return(returnArgs[0]...)
			authN.On("FetchAuthFromToken", token).Return(returnArgs[1]...)

			authDirective := authdirective.New(db, authN, 0)

			_, err := authDirective.AuthMiddleware()(ctx, nil, next)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}

func (ts *TestSuite) TestAuthRenewalMiddleware() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	ctx := context.Background()

	tokenString := ""

	token := &jwt.Token{}
	timeBeforeTokenExpTimeInSec := 0

	next := func(ctx context.Context) (interface{}, error) { return nil, nil }

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithAuthRenewalMiddleware",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     id,
					"userID": userID,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				authDatastore := datastoremodelfactory.NewAuth(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheAuthorizationTokenIsNotSetInTheContext",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				returnArgs = ReturnArgs{
					{nil, nil},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheTokenIsNotValidForRenewal",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				returnArgs = ReturnArgs{
					{nil, customerror.New("failed")},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFetchedFromTheToken",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodel.Auth{}, customerror.New("failed")},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenTryingToFindTheAuthInTheDatabase",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnError(customerror.New("failed"))
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFoundInTheDatabase",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(&sqlmock.Rows{})
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheUserIDFromTokenDoesNotMatchWithTheUserIDFromAuthRecordFromTheDatabase",
			SetUp: func(t *testing.T) {
				ctx = context.Background()

				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				authDatastore := datastoremodelfactory.NewAuth(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authN := new(mockauthpkg.Auth)
			authN.On("ValidateTokenRenewal", tokenString, timeBeforeTokenExpTimeInSec).Return(returnArgs[0]...)
			authN.On("FetchAuthFromToken", token).Return(returnArgs[1]...)

			authDirective := authdirective.New(db, authN, timeBeforeTokenExpTimeInSec)

			_, err := authDirective.AuthRenewalMiddleware()(ctx, nil, next)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}
