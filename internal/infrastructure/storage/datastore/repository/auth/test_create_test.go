package auth_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	datastoreentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestCreate() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	auth := domainmodel.Auth{}

	newAuth := domainmodel.Auth{}

	errorType := customerror.NoType

	firstSqlQuery := `INSERT INTO "auths" ("user_id","created_at","id") VALUES ($1,$2,$3) RETURNING "id"`

	secondSqlQuery := `SELECT * FROM "logins" WHERE user_id=$1`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingAnAuth",
			SetUp: func(t *testing.T) {
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     uuid.Nil,
					"userID": userID,
				}

				auth = domainmodelfactory.NewAuth(args)

				args = map[string]interface{}{
					"id":     uuid.Nil,
					"userID": auth.UserID,
				}

				newAuth = domainmodelfactory.NewAuth(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(firstSqlQuery)).
					WithArgs(auth.UserID, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAnAuth",
			SetUp: func(t *testing.T) {
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     uuid.Nil,
					"userID": userID,
				}

				auth = domainmodelfactory.NewAuth(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(firstSqlQuery)).
					WithArgs(auth.UserID, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAnAuthBecauseTheUserAuthIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     uuid.Nil,
					"userID": userID,
				}

				auth = domainmodelfactory.NewAuth(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(firstSqlQuery)).
					WithArgs(auth.UserID, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.Conflict.New("auths_user_id_key"))

				mock.ExpectRollback()

				args = map[string]interface{}{
					"userID": userID,
				}

				login := datastoreentityfactory.NewLogin(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(login.ID, login.UserID, login.Username, login.Password, login.CreatedAt, login.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(secondSqlQuery)).
					WithArgs(auth.UserID).
					WillReturnRows(rows)

				errorType = customerror.Conflict
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAnAuthBecauseTheUserAuthIsAlreadyRegisteredAndLoginIsNotFound",
			SetUp: func(t *testing.T) {
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     uuid.Nil,
					"userID": userID,
				}

				auth = domainmodelfactory.NewAuth(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(firstSqlQuery)).
					WithArgs(auth.UserID, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("auths_user_id_key"))

				mock.ExpectRollback()

				mock.ExpectQuery(regexp.QuoteMeta(secondSqlQuery)).
					WithArgs(auth.UserID).
					WillReturnRows(&sqlmock.Rows{})

				errorType = customerror.NotFound
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAnAuthBecauseTheUserAuthIsAlreadyRegisteredAndAnErrorAlsoHappensWhenFindingTheLoginByUserID",
			SetUp: func(t *testing.T) {
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     uuid.Nil,
					"userID": userID,
				}

				auth = domainmodelfactory.NewAuth(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(firstSqlQuery)).
					WithArgs(auth.UserID, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("auths_user_id_key"))

				mock.ExpectRollback()

				mock.ExpectQuery(regexp.QuoteMeta(secondSqlQuery)).
					WithArgs(auth.UserID).
					WillReturnError(customerror.New("failed"))

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(db)

			returnedAuth, err := authDatastoreRepository.Create(auth)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, newAuth.UserID, returnedAuth.UserID)
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
