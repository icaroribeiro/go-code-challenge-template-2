package user_test

import (
	"testing"

	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/user"
	userdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/mockrepository/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockvalidator"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestWithDBTrx() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrx := &gorm.DB{}

	userDatastoreRepositoryWithDBTrx := &userdatastoremockrepository.Repository{}

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSettingRepositoriesWithDatabaseTransaction",
			SetUp: func(t *testing.T) {
				dbTrx = db

				userDatastoreRepositoryWithDBTrx = &userdatastoremockrepository.Repository{}

				returnArgs = ReturnArgs{
					{userDatastoreRepositoryWithDBTrx},
				}
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)
			userDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[0]...)

			validator := new(mockvalidator.Validator)

			userService := userservice.New(userDatastoreRepository, validator)

			returnedUserService := userService.WithDBTrx(dbTrx)

			if !tc.WantError {
				assert.NotEmpty(t, returnedUserService, "Service interface is empty.")
				assert.Equal(t, userService, returnedUserService, "Service interfaces are not the same.")
			}
		})
	}
}
