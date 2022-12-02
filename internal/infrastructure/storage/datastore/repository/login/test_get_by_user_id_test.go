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
