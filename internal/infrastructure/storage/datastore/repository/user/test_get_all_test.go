package user_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/entity"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestGetAll() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	user := domainentity.User{}

	errorType := customerror.NoType

	sqlQuery := `SELECT * FROM "users"`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				userDatastore := datastoremodelfactory.NewUser(nil)
				user = userDatastore.ToDomain()

				rows := sqlmock.
					NewRows([]string{"id", "username", "created_at", "updated_at"}).
					AddRow(userDatastore.ID, userDatastore.Username, userDatastore.CreatedAt, userDatastore.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenFindingAllUser",
			SetUp: func(t *testing.T) {
				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WillReturnError(customerror.New("failed"))

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			userDatastoreRepository := userdatastorerepository.New(db)

			returnedUsers, err := userDatastoreRepository.GetAll()

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, user.ID, returnedUsers[0].ID)
				assert.Equal(t, user.Username, returnedUsers[0].Username)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedUsers)
			}
		})
	}
}
