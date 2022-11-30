package datastore_test

import (
	"fmt"
	"testing"

	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/datastore"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNew() {
	dbConfig := map[string]string{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInInitializingThePostgresDriver",
			SetUp: func(t *testing.T) {
				dbConfig = ts.PostgresDBConfig
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheSQLDatabaseDriverIsNotRecognized",
			SetUp: func(t *testing.T) {
				dbConfig["DRIVER"] = "testing"
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			_, err := datastorepkg.New(dbConfig)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}
