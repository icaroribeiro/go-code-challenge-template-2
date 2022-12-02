package auth_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/auth"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	authdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/mockrepository/auth"
	logindatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/mockrepository/login"
	userdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/mockrepository/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
	mockauth "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockauth"
	mocksecurity "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mocksecurity"
	mockvalidator "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockvalidator"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestRegister() {
	credentials := security.Credentials{}

	user := domainentity.User{}

	login := domainentity.Login{}

	auth := domainentity.Auth{}

	newAuth := domainentity.Auth{}

	tokenExpTimeInSec := fake.Number(2, 10)

	token := ""

	errorType := customerror.NoType

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRegisteringAUser",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				args := map[string]interface{}{
					"id":       uuid.Nil,
					"username": credentials.Username,
				}

				user = domainentityfactory.NewUser(args)

				id := uuid.NewV4()

				args = map[string]interface{}{
					"id":       id,
					"username": credentials.Username,
				}

				newUser := domainentityfactory.NewUser(args)

				args = map[string]interface{}{
					"id":       uuid.Nil,
					"userID":   newUser.ID,
					"username": credentials.Username,
					"password": credentials.Password,
				}

				login = domainentityfactory.NewLogin(args)

				args = map[string]interface{}{
					"id":     uuid.Nil,
					"userID": newUser.ID,
				}

				auth = domainentityfactory.NewAuth(args)

				id = uuid.NewV4()

				args = map[string]interface{}{
					"id":     id,
					"userID": newUser.ID,
				}

				newAuth = domainentityfactory.NewAuth(args)

				token = fake.Word()

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Login{}, nil},
					{newUser, nil},
					{domainentity.Login{}, nil},
					{newAuth, nil},
					{token, nil},
				}
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheCredentialsAreNotValid",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
					{domainentity.Login{}, nil},
					{domainentity.User{}, nil},
					{domainentity.Login{}, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingALoginByUsername",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = security.Credentials{
					Username: username,
					Password: password,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Login{}, customerror.New("failed")},
					{domainentity.User{}, nil},
					{domainentity.Login{}, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheUsernameIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainentity.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{domainentity.User{}, nil},
					{domainentity.Login{}, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.Conflict
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAUser",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				user = domainentity.User{
					Username: credentials.Username,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Login{}, nil},
					{domainentity.User{}, customerror.New("failed")},
					{domainentity.Login{}, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingALogin",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				user = domainentity.User{
					Username: credentials.Username,
				}

				id := uuid.NewV4()

				newUser := domainentity.User{
					ID:       id,
					Username: credentials.Username,
				}

				login = domainentity.Login{
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Login{}, nil},
					{newUser, nil},
					{domainentity.Login{}, customerror.New("failed")},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAnAuth",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				user = domainentity.User{
					Username: credentials.Username,
				}

				id := uuid.NewV4()

				newUser := domainentity.User{
					ID:       id,
					Username: credentials.Username,
				}

				login = domainentity.Login{
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				id = uuid.NewV4()

				newLogin := domainentity.Login{
					ID:       id,
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainentity.Auth{
					UserID: newUser.ID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Login{}, nil},
					{newUser, nil},
					{newLogin, nil},
					{domainentity.Auth{}, customerror.New("failed")},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAToken",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				user = domainentity.User{
					Username: credentials.Username,
				}

				id := uuid.NewV4()

				newUser := domainentity.User{
					ID:       id,
					Username: credentials.Username,
				}

				login = domainentity.Login{
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				id = uuid.NewV4()

				newLogin := domainentity.Login{
					ID:       id,
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainentity.Auth{
					UserID: newUser.ID,
				}

				id = uuid.NewV4()

				newAuth = domainentity.Auth{
					ID:     id,
					UserID: newUser.ID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Login{}, nil},
					{newUser, nil},
					{newLogin, nil},
					{newAuth, nil},
					{"", customerror.New("failed")},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			validator := new(mockvalidator.Validator)
			validator.On("Validate", credentials).Return(returnArgs[0]...)

			loginDatastoreRepository := new(logindatastoremockrepository.Repository)
			loginDatastoreRepository.On("GetByUsername", credentials.Username).Return(returnArgs[1]...)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)
			userDatastoreRepository.On("Create", user).Return(returnArgs[2]...)

			loginDatastoreRepository.On("Create", login).Return(returnArgs[3]...)

			authDatastoreRepository := new(authdatastoremockrepository.Repository)
			authDatastoreRepository.On("Create", auth).Return(returnArgs[4]...)

			authN := new(mockauth.Auth)
			authN.On("CreateToken", newAuth, tokenExpTimeInSec).Return(returnArgs[5]...)

			security := new(mocksecurity.Security)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, security, validator, tokenExpTimeInSec)

			returnedToken, err := authService.Register(credentials)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, token, returnedToken)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedToken)
			}

			tc.TearDown(t)
		})
	}
}
