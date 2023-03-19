package login_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	persistententity "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/perentity"
	logindatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/repository/login"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

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

				persistentLogin := persistententity.LoginFactory(args)
				login = persistentLogin.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(persistentLogin.ID, persistentLogin.UserID, persistentLogin.Username, persistentLogin.Password, persistentLogin.CreatedAt, persistentLogin.UpdatedAt)

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

				persistentLogin := persistententity.LoginFactory(args)
				login = persistentLogin.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(persistentLogin.ID, persistentLogin.UserID, persistentLogin.Username, persistentLogin.Password, persistentLogin.CreatedAt, persistentLogin.UpdatedAt)

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

				persistentLogin := persistententity.LoginFactory(args)
				login = persistentLogin.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
					AddRow(persistentLogin.ID, persistentLogin.UserID, persistentLogin.Username, persistentLogin.Password, persistentLogin.CreatedAt, persistentLogin.UpdatedAt)

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

			persistentLoginRepository := logindatastorerepository.New(db)

			returnedLogin, err := persistentLoginRepository.Delete(id.String())

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
