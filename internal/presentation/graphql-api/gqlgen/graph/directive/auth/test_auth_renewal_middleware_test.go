package auth_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	authdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/directive/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/entity"
	mockauthpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockauth"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestAuthRenewalMiddleware() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	ctx := context.Background()

	token := &jwt.Token{}
	timeBeforeTokenExpTimeInSec := 0

	next := func(ctx context.Context) (interface{}, error) { return nil, nil }

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithAuthRenewalMiddleware",
			SetUp: func(t *testing.T) {
				ctx = context.Background()
				ctx = authmiddlewarepkg.NewContext(ctx, token)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     id,
					"userID": userID,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainentityfactory.NewAuth(args), nil},
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
					{domainentity.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheTokenIsNotValidForRenewal",
			SetUp: func(t *testing.T) {
				ctx = context.Background()
				ctx = authmiddlewarepkg.NewContext(ctx, token)

				returnArgs = ReturnArgs{
					{nil, customerror.New("failed")},
					{domainentity.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFetchedFromTheToken",
			SetUp: func(t *testing.T) {
				ctx = context.Background()
				ctx = authmiddlewarepkg.NewContext(ctx, token)

				returnArgs = ReturnArgs{
					{token, nil},
					{domainentity.Auth{}, customerror.New("failed")},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenTryingToFindTheAuthInTheDatabase",
			SetUp: func(t *testing.T) {
				ctx = context.Background()
				ctx = authmiddlewarepkg.NewContext(ctx, token)

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainentityfactory.NewAuth(args), nil},
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
				ctx = authmiddlewarepkg.NewContext(ctx, token)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainentityfactory.NewAuth(args), nil},
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
				ctx = authmiddlewarepkg.NewContext(ctx, token)

				token = &jwt.Token{}

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainentityfactory.NewAuth(args), nil},
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
			authN.On("ValidateTokenRenewal", token, timeBeforeTokenExpTimeInSec).Return(returnArgs[0]...)
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
