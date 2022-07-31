package login_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	fake "github.com/brianvoe/gofakeit/v5"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/login"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/entity"
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

	login := domainentity.Login{}

	newLogin := domainentity.Login{}

	errorType := customerror.NoType

	stmt := `INSERT INTO "logins" ("id","user_id","username","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6)`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingTheLogin",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainentityfactory.NewLogin(args)

				args = map[string]interface{}{
					"userID":   login.UserID,
					"username": login.Username,
					"password": login.Password,
				}

				newLogin = domainentityfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(stmt)).
					WithArgs(sqlmock.AnyArg(), login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheLogin",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainentityfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(stmt)).
					WithArgs(sqlmock.AnyArg(), login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheLoginBecauseTheUserLoginIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainentityfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(stmt)).
					WithArgs(sqlmock.AnyArg(), login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.Conflict.New("logins_user_id_key"))

				mock.ExpectRollback()

				errorType = customerror.Conflict
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			loginDatastoreRepository := logindatastorerepository.New(db)

			returnedLogin, err := loginDatastoreRepository.Create(login)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, newLogin.UserID, returnedLogin.UserID)
				assert.Equal(t, newLogin.Username, returnedLogin.Username)
				security := securitypkg.New()
				err := security.VerifyPasswords(returnedLogin.Password, newLogin.Password)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedLogin)
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}

func (ts *TestSuite) TestGetByUsername() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	username := ""

	login := domainentity.Login{}

	errorType := customerror.NoType

	stmt := `SELECT * FROM "logins" WHERE username=$1`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingTheLoginByUsername",
			SetUp: func(t *testing.T) {
				username = fake.Username()

				args := map[string]interface{}{
					"username": username,
				}

				loginDatastore := datastoremodelfactory.NewLogin(args)
				login = loginDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(loginDatastore.ID, loginDatastore.UserID, loginDatastore.Username, loginDatastore.Password, loginDatastore.CreatedAt, loginDatastore.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(stmt)).
					WithArgs(username).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldSucceedIfTheUsernameIsNotFound",
			SetUp: func(t *testing.T) {
				username = fake.Username()

				args := map[string]interface{}{
					"id":       uuid.Nil,
					"userID":   uuid.Nil,
					"username": "",
					"password": "",
				}

				login = domainentityfactory.NewLogin(args)

				mock.ExpectQuery(regexp.QuoteMeta(stmt)).
					WithArgs(username).
					WillReturnRows(&sqlmock.Rows{})
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheLoginByUsername",
			SetUp: func(t *testing.T) {
				username = fake.Username()

				mock.ExpectQuery(regexp.QuoteMeta(stmt)).
					WithArgs(username).
					WillReturnError(customerror.New("failed"))

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			loginDatastoreRepository := logindatastorerepository.New(db)

			returnedLogin, err := loginDatastoreRepository.GetByUsername(username)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, login.ID, returnedLogin.ID)
				assert.Equal(t, login.UserID, returnedLogin.UserID)
				assert.Equal(t, login.Username, returnedLogin.Username)
				assert.Equal(t, login.Password, returnedLogin.Password)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedLogin)
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

	login := domainentity.Login{}

	errorType := customerror.NoType

	stmt := `SELECT * FROM "logins" WHERE user_id=$1`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingTheLoginByUserID",
			SetUp: func(t *testing.T) {
				userID = uuid.NewV4()

				args := map[string]interface{}{
					"userID": userID,
				}

				loginDatastore := datastoremodelfactory.NewLogin(args)
				login = loginDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(loginDatastore.ID, loginDatastore.UserID, loginDatastore.Username, loginDatastore.Password, loginDatastore.CreatedAt, loginDatastore.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(stmt)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldSucceedIfTheLoginIsNotFound",
			SetUp: func(t *testing.T) {
				userID = uuid.NewV4()

				args := map[string]interface{}{
					"id":       uuid.Nil,
					"userID":   uuid.Nil,
					"username": "",
					"password": "",
				}

				login = domainentityfactory.NewLogin(args)

				mock.ExpectQuery(regexp.QuoteMeta(stmt)).
					WithArgs(userID).
					WillReturnRows(&sqlmock.Rows{})
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheLoginByUserID",
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

			loginDatastoreRepository := logindatastorerepository.New(db)

			returnedLogin, err := loginDatastoreRepository.GetByUserID(userID.String())

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, login.ID, returnedLogin.ID)
				assert.Equal(t, login.UserID, returnedLogin.UserID)
				assert.Equal(t, login.Username, returnedLogin.Username)
				assert.Equal(t, login.Password, returnedLogin.Password)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedLogin)
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}

func (ts *TestSuite) TestUpdate() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	var id uuid.UUID

	login := domainentity.Login{}

	updatedLogin := domainentity.Login{}

	errorType := customerror.NoType

	firstStmt := `UPDATE "logins" SET "user_id"=$1,"username"=$2,"password"=$3,"updated_at"=$4 WHERE id=$5`

	secondStmt := `SELECT * FROM "logins" WHERE id=$1`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInUpdatingTheLogin",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainentityfactory.NewLogin(args)

				args = map[string]interface{}{
					"id":     id,
					"userID": login.UserID,
				}

				updatedDatastoreLogin := datastoremodelfactory.NewLogin(args)
				updatedLogin = updatedDatastoreLogin.ToDomain()

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(id, updatedDatastoreLogin.UserID, updatedDatastoreLogin.Username, updatedDatastoreLogin.Password, updatedDatastoreLogin.CreatedAt, updatedDatastoreLogin.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(secondStmt)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenUpdatingTheLogin",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainentityfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheLoginIsNotFound",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainentityfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()

				errorType = customerror.NotFound
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheLoginByID",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainentityfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()

				mock.ExpectQuery(regexp.QuoteMeta(secondStmt)).
					WithArgs(id).
					WillReturnError(customerror.New("failed"))

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheLoginIsNotFoundAfterUpdatingIt",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainentityfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(firstStmt)).
					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()

				mock.ExpectQuery(regexp.QuoteMeta(secondStmt)).
					WithArgs(id).
					WillReturnRows(&sqlmock.Rows{})

				errorType = customerror.NotFound
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			loginDatastoreRepository := logindatastorerepository.New(db)

			returnedLogin, err := loginDatastoreRepository.Update(id.String(), login)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, updatedLogin.ID, returnedLogin.ID)
				assert.Equal(t, updatedLogin.UserID, returnedLogin.UserID)
				assert.Equal(t, updatedLogin.Username, returnedLogin.Username)
				assert.Equal(t, updatedLogin.Password, returnedLogin.Password)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedLogin)
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

	login := domainentity.Login{}

	errorType := customerror.NoType

	firstStmt := `SELECT * FROM "logins" WHERE id=$1`

	secondStmt := `DELETE FROM "logins" WHERE "logins"."id" = $1`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInDeletingTheLogin",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				loginDatastore := datastoremodelfactory.NewLogin(args)
				login = loginDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(loginDatastore.ID, loginDatastore.UserID, loginDatastore.Username, loginDatastore.Password, loginDatastore.CreatedAt, loginDatastore.UpdatedAt)

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
			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheLoginByID",
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
			Context: "ItShouldFailIfTheLoginIsNotFound",
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
			Context: "ItShouldFailIfAnErrorOccursWhenDeletingTheLoginByID",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				loginDatastore := datastoremodelfactory.NewLogin(args)
				login = loginDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(loginDatastore.ID, loginDatastore.UserID, loginDatastore.Username, loginDatastore.Password, loginDatastore.CreatedAt, loginDatastore.UpdatedAt)

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
			Context: "ItShouldFailIfTheLoginIsNotDeleted",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				loginDatastore := datastoremodelfactory.NewLogin(args)
				login = loginDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(loginDatastore.ID, loginDatastore.UserID, loginDatastore.Username, loginDatastore.Password, loginDatastore.CreatedAt, loginDatastore.UpdatedAt)

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

			loginDatastoreRepository := logindatastorerepository.New(db)

			returnedLogin, err := loginDatastoreRepository.Delete(id.String())

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, login.ID, returnedLogin.ID)
				assert.Equal(t, login.UserID, returnedLogin.UserID)
				assert.Equal(t, login.Username, returnedLogin.Username)
				assert.Equal(t, login.Password, returnedLogin.Password)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedLogin)
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

			loginDatastoreRepository := logindatastorerepository.New(db)

			returnedLoginDatastoreRepository := loginDatastoreRepository.WithDBTrx(dbTrx)

			if !tc.WantError {
				assert.NotEmpty(t, returnedLoginDatastoreRepository, "Repository interface is empty.")
				assert.Equal(t, loginDatastoreRepository, returnedLoginDatastoreRepository, "Repository interfaces are not the same.")
			}
		})
	}
}
