package user_test

import (
	"testing"

	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/user"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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
