package datastore_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestPostgresDriver(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestNewPostgresDriver() {
	dbConfig := map[string]string{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingAPostgresDriverWithDBConfig",
			SetUp: func(t *testing.T) {
				dbConfig = ts.PostgresDBConfig
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldSucceedInCreatingAPostgresDriverUsingDBConfigURL",
			SetUp: func(t *testing.T) {
				dbConfig = map[string]string{
					"URL": ts.PostgresDBConfigURL,
				}
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailInCreatingAPostgresDriverUsingDBConfigURLIfTheURLIsInvalid",
			SetUp: func(t *testing.T) {
				dbConfig = map[string]string{
					"URL": "testing",
				}
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			driver, err := datastorepkg.NewPostgresDriver(dbConfig)

			if !tc.WantError {
				assert.NotEmpty(t, driver)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.Empty(t, driver)
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}

}

func (ts *TestSuite) TestGetInstancePostgresDriver() {
	dbConfig := map[string]string{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAPostgresDBInstance",
			SetUp: func(t *testing.T) {
				dbConfig = ts.PostgresDBConfig
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailInGettingAPostgresDBInstanceIfTheURLIsInvalid",
			SetUp: func(t *testing.T) {
				dbConfig = map[string]string{
					"URL": "testing",
				}
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			driver, err := datastorepkg.NewPostgresDriver(dbConfig)

			if !tc.WantError {
				assert.NotEmpty(t, driver)
				db := driver.GetInstance()
				assert.NotNil(t, db, "Database is nil")
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.Empty(t, driver)
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}

func (ts *TestSuite) TestClosePostgresDriver() {
	driver := "postgres"
	db := &gorm.DB{}
	var mock sqlmock.Sqlmock
	provider := datastorepkg.Provider{}

	var connPool gorm.ConnPool

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInClosingTheDatabase",
			SetUp: func(t *testing.T) {
				db, mock = NewMockDB(driver)
				provider = datastorepkg.Provider{DB: db}
				mock.ExpectClose()
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingTheSQLDatabase",
			SetUp: func(t *testing.T) {
				db, _ = NewMockDB(driver)
				connPool = db.ConnPool
				db.ConnPool = nil
				provider = datastorepkg.Provider{DB: db}
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				db.ConnPool = connPool
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			postgresDriver := datastorepkg.PostgresDriver{Provider: provider}

			err := postgresDriver.Close()

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))

			tc.TearDown(t)
		})
	}
}
