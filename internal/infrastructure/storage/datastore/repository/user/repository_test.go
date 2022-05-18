package user_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainfactorymodel "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
	datastorefactorymodel "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/model"
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

	user := domainmodel.User{}

	newUser := domainmodel.User{}

	errorType := customerror.NoType

	sqlQuery := `INSERT INTO "users" ("id","username","created_at","updated_at") VALUES ($1,$2,$3,$4)`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingTheUser",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				user = domainfactorymodel.NewUser(args)

				args = map[string]interface{}{
					"id":       uuid.Nil,
					"username": user.Username,
				}

				newUser = domainfactorymodel.NewUser(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

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

				user = domainfactorymodel.NewUser(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheUserBecauseTheUserIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				user = domainfactorymodel.NewUser(nil)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			userDatastoreRepository := userdatastorerepository.New(db)

			returnedUser, err := userDatastoreRepository.Create(user)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
				assert.Equal(t, newUser.Username, returnedUser.Username)
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

func (ts *TestSuite) TestGetAll() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	user := domainmodel.User{}

	errorType := customerror.NoType

	sqlQuery := `SELECT * FROM "users"`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				userDatastore := datastorefactorymodel.NewUser(nil)
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
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
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

			userDatastoreRepository := userdatastorerepository.New(db)

			returnedUserDatastoreRepository := userDatastoreRepository.WithDBTrx(dbTrx)

			if !tc.WantError {
				assert.NotEmpty(t, returnedUserDatastoreRepository, "Repository interface is empty.")
				assert.Equal(t, userDatastoreRepository, returnedUserDatastoreRepository, "Repository interfaces are not the same.")
			}
		})
	}
}
