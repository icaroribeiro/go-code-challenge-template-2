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

func (ts *TestSuite) TestGetByUserID() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	var userID uuid.UUID

	auth := domainentity.Auth{}

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
