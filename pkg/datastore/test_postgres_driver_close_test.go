package datastore_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/datastore"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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
