package user_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/repository/user"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestCreate() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	user := domainentity.User{}

	UserFactory := domainentity.User{}

	errorType := customerror.NoType

	sqlQuery := `INSERT INTO "users" ("username","created_at","updated_at","id") VALUES ($1,$2,$3,$4) RETURNING "id"`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingTheUser",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				user = domainentity.UserFactory(args)

				args = map[string]interface{}{
					"id":       uuid.Nil,
					"username": user.Username,
				}

				UserFactory = domainentity.UserFactory(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheUser",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				user = domainentity.UserFactory(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheUserBecauseTheUserIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				user = domainentity.UserFactory(nil)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.Conflict.New("duplicate key value"))

				mock.ExpectRollback()

				errorType = customerror.Conflict
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			persistentUserRepository := userdatastorerepository.New(db)

			returnedUser, err := persistentUserRepository.Create(user)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, UserFactory.Username, returnedUser.Username)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedUser)
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}
