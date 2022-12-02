package auth_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestDelete() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	var id uuid.UUID

	auth := domainentity.Auth{}

	errorType := customerror.NoType

	firstStmt := `SELECT * FROM "auths" WHERE id=$1`

	secondStmt := `DELETE FROM "auths" WHERE "auths"."id" = $1`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInDeletingAnAuth",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				authDatastore := datastoremodelfactory.NewAuth(args)
				auth = authDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(firstStmt)).
					WithArgs(id).
					WillReturnRows(rows)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(secondStmt)).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenFindingAnAuthByID",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				mock.ExpectQuery(regexp.QuoteMeta(firstStmt)).
					WithArgs(id).
					WillReturnError(customerror.New("failed"))

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnAuthIsNotFound",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				mock.ExpectQuery(regexp.QuoteMeta(firstStmt)).
					WithArgs(id).
					WillReturnRows(&sqlmock.Rows{})

				errorType = customerror.NotFound
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenDeletingAnAuthByID",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				authDatastore := datastoremodelfactory.NewAuth(args)
				auth = authDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(firstStmt)).
					WithArgs(id).
					WillReturnRows(rows)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(secondStmt)).
					WithArgs(id).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnAuthIsNotDeleted",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				authDatastore := datastoremodelfactory.NewAuth(args)
				auth = authDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(firstStmt)).
					WithArgs(id).
					WillReturnRows(rows)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(secondStmt)).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()

				errorType = customerror.NotFound
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(db)

			returnedAuth, err := authDatastoreRepository.Delete(id.String())

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, auth.ID, returnedAuth.ID)
				assert.Equal(t, auth.UserID, returnedAuth.UserID)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedAuth)
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}
