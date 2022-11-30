package datastore_test

import (
	"fmt"
	"testing"

	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/datastore"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestPostgresDriverGetInstance() {
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
