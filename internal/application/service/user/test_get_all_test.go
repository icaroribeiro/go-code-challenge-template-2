package user_test

import (
	"fmt"
	"testing"

	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/user"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	userdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/mockrepository/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockvalidator"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestGetAll() {
	user := domainentity.User{}

	returnArgs := ReturnArgs{}

	errorType := customerror.NoType

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				user = domainentityfactory.NewUser(nil)

				returnArgs = ReturnArgs{
					{domainentity.Users{user}, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetAllUsers",
			SetUp: func(t *testing.T) {
				returnArgs = ReturnArgs{
					{domainentity.Users{}, customerror.New("failed")},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)
			userDatastoreRepository.On("GetAll").Return(returnArgs[0]...)

			validator := new(mockvalidator.Validator)

			userService := userservice.New(userDatastoreRepository, validator)

			returnedUsers, err := userService.GetAll()

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
