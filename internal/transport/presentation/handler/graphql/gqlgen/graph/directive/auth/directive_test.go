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
	authdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	domainfactorymodel "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
	datastorefactorymodel "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/model"
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
				authDetailsCtxValue = domainfactorymodel.NewAuth(nil)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedCtx := authdirectivepkg.NewContext(ctx, authDetailsCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedAuthDetailsCtxValue, ok := authdirectivepkg.FromContext(returnedCtx)
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
			Context: "ItShouldSucceedInGettingAssociatedValueWithAContext",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainfactorymodel.NewAuth(nil)
				ctx = authdirectivepkg.NewContext(ctx, authDetailsCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedAuthDetailsCtxValue, ok := authdirectivepkg.FromContext(ctx)

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

	tokenString := ""

	ctx := context.Background()

	token := &jwt.Token{}

	next := func(ctx context.Context) (interface{}, error) { return nil, nil }

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithAuthenticationMiddleware",
			SetUp: func(t *testing.T) {
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
					{domainfactorymodel.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				authDatastore := datastorefactorymodel.NewAuth(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: false,
			TearDown: func(t *testing.T) {
				ctx = context.Background()
			},
		},
		{
			Context: "ItShouldFailIfTheAuthenticationTokenIsNotSetInTheContext",
			SetUp: func(t *testing.T) {
				returnArgs = ReturnArgs{
					{nil, nil},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheTokenIsNotDecoded",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				returnArgs = ReturnArgs{
					{nil, customerror.New("failed")},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				ctx = context.Background()
			},
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFetchedFromTheToken",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodel.Auth{}, customerror.New("failed")},
				}
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				ctx = context.Background()
			},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenTryingToFindTheAuthInTheDatabase",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainfactorymodel.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnError(customerror.New("failed"))
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				ctx = context.Background()
			},
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFoundInTheDatabase",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainfactorymodel.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(&sqlmock.Rows{})
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				ctx = context.Background()
			},
		},
		{
			Context: "ItShouldFailIfTheUserIDFromTokenDoesNotMatchWithTheUserIDFromAuthRecordFromTheDatabase",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()
				ctx = authmiddlewarepkg.NewContext(ctx, tokenString)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainfactorymodel.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				authDatastore := datastorefactorymodel.NewAuth(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				ctx = context.Background()
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authN := new(mockauthpkg.Auth)
			authN.On("DecodeToken", tokenString).Return(returnArgs[0]...)
			authN.On("FetchAuthFromToken", token).Return(returnArgs[1]...)

			authMiddleware := authdirectivepkg.AuthMiddleware(db, authN)

			_, err := authMiddleware(ctx, nil, next)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))

			tc.TearDown(t)
		})
	}
}

// func (ts *TestSuite) TestAuthRenewalMiddleware() {}
