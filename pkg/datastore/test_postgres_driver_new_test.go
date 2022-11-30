package datastore_test

import (
	"fmt"
	"testing"

	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/datastore"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestPostgresDriverNew() {
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
