package auth_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestRepository(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestCreate() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	auth := domainmodel.Auth{}

	newAuth := domainmodel.Auth{}

	errorType := customerror.NoType

	firstStmt := `INSERT INTO "auths" ("id","user_id","created_at") VALUES ($1,$2,$3)`

	secondStmt := `SELECT * FROM "logins" WHERE user_id=$1`

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

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

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

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
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

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
					WillReturnError(customerror.Conflict.New("auths_user_id_key"))

				mock.ExpectRollback()

				args = map[string]interface{}{
					"userID": userID,
				}

				login := datastoremodelfactory.NewLogin(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(login.ID, login.UserID, login.Username, login.Password, login.CreatedAt, login.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(secondStmt)).
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

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
					WillReturnError(customerror.New("auths_user_id_key"))

				mock.ExpectRollback()

				mock.ExpectQuery(regexp.QuoteMeta(secondStmt)).
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

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
					WillReturnError(customerror.New("auths_user_id_key"))

				mock.ExpectRollback()

				mock.ExpectQuery(regexp.QuoteMeta(secondStmt)).
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

func (ts *TestSuite) TestGetByUserID() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	var userID uuid.UUID

	auth := domainmodel.Auth{}

	errorType := customerror.NoType

	stmt := `SELECT * FROM "auths" WHERE user_id=$1`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAnAuthByUserID",
			SetUp: func(t *testing.T) {
				userID = uuid.NewV4()

				args := map[string]interface{}{
					"userID": userID,
				}

				authDatastore := datastoremodelfactory.NewAuth(args)
				auth = authDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(stmt)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenFindingAnAuthByUserID",
			SetUp: func(t *testing.T) {
				userID = uuid.NewV4()

				mock.ExpectQuery(regexp.QuoteMeta(stmt)).
					WithArgs(userID).
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

			returnedAuth, err := authDatastoreRepository.GetByUserID(userID.String())

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

func (ts *TestSuite) TestDelete() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	var id uuid.UUID

	auth := domainmodel.Auth{}

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

func (ts *TestSuite) TestWithDBTrx() {
	driver := "postgres"
	db, _ := NewMockDB(driver)

	dbTrx := &gorm.DB{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSettingTheRepositoryWithDatabaseTransaction",
			SetUp: func(t *testing.T) {
				dbTrx = db.Begin()
			},
			WantError: false,
		},
		{
			Context: "ItShouldSucceedInSettingTheRepositoryWithoutDatabaseTransaction",
			SetUp: func(t *testing.T) {
				dbTrx = nil
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(db)

			returnedAuthDatastoreRepository := authDatastoreRepository.WithDBTrx(dbTrx)

			if !tc.WantError {
				assert.NotEmpty(t, returnedAuthDatastoreRepository, "Repository interface is empty.")
				assert.Equal(t, authDatastoreRepository, returnedAuthDatastoreRepository, "Repository interfaces are not the same.")
			}
		})
	}
}
