package login_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/login"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

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
