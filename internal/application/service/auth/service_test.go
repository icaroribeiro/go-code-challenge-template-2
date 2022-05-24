package auth_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/auth"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/mockrepository/auth"
	logindatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/mockrepository/login"
	userdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/mockrepository/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
	domainfactorymodel "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
	mockauth "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockauth"
	mocksecurity "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mocksecurity"
	mockvalidator "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockvalidator"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestServiceUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestRegister() {
	credentials := security.Credentials{}

	user := domainmodel.User{}

	login := domainmodel.Login{}

	auth := domainmodel.Auth{}

	newAuth := domainmodel.Auth{}

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

				user = domainfactorymodel.NewUser(args)

				id := uuid.NewV4()

				args = map[string]interface{}{
					"id":       id,
					"username": credentials.Username,
				}

				newUser := domainfactorymodel.NewUser(args)

				args = map[string]interface{}{
					"id":       uuid.Nil,
					"userID":   newUser.ID,
					"username": credentials.Username,
					"password": credentials.Password,
				}

				login = domainfactorymodel.NewLogin(args)

				args = map[string]interface{}{
					"id":     uuid.Nil,
					"userID": newUser.ID,
				}

				auth = domainfactorymodel.NewAuth(args)

				id = uuid.NewV4()

				args = map[string]interface{}{
					"id":     id,
					"userID": newUser.ID,
				}

				newAuth = domainfactorymodel.NewAuth(args)

				token = fake.Word()

				returnArgs = ReturnArgs{
					{nil},
					{domainmodel.Login{}, nil},
					{newUser, nil},
					{domainmodel.Login{}, nil},
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
					{domainmodel.Login{}, nil},
					{domainmodel.User{}, nil},
					{domainmodel.Login{}, nil},
					{domainmodel.Auth{}, nil},
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
					{domainmodel.Login{}, customerror.New("failed")},
					{domainmodel.User{}, nil},
					{domainmodel.Login{}, nil},
					{domainmodel.Auth{}, nil},
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

				login = domainmodel.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{domainmodel.User{}, nil},
					{domainmodel.Login{}, nil},
					{domainmodel.Auth{}, nil},
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

				user = domainmodel.User{
					Username: credentials.Username,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainmodel.Login{}, nil},
					{domainmodel.User{}, customerror.New("failed")},
					{domainmodel.Login{}, nil},
					{domainmodel.Auth{}, nil},
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

				user = domainmodel.User{
					Username: credentials.Username,
				}

				id := uuid.NewV4()

				newUser := domainmodel.User{
					ID:       id,
					Username: credentials.Username,
				}

				login = domainmodel.Login{
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainmodel.Login{}, nil},
					{newUser, nil},
					{domainmodel.Login{}, customerror.New("failed")},
					{domainmodel.Auth{}, nil},
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

				user = domainmodel.User{
					Username: credentials.Username,
				}

				id := uuid.NewV4()

				newUser := domainmodel.User{
					ID:       id,
					Username: credentials.Username,
				}

				login = domainmodel.Login{
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				id = uuid.NewV4()

				newLogin := domainmodel.Login{
					ID:       id,
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainmodel.Auth{
					UserID: newUser.ID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainmodel.Login{}, nil},
					{newUser, nil},
					{newLogin, nil},
					{domainmodel.Auth{}, customerror.New("failed")},
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

				user = domainmodel.User{
					Username: credentials.Username,
				}

				id := uuid.NewV4()

				newUser := domainmodel.User{
					ID:       id,
					Username: credentials.Username,
				}

				login = domainmodel.Login{
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				id = uuid.NewV4()

				newLogin := domainmodel.Login{
					ID:       id,
					UserID:   newUser.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainmodel.Auth{
					UserID: newUser.ID,
				}

				id = uuid.NewV4()

				newAuth = domainmodel.Auth{
					ID:     id,
					UserID: newUser.ID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{domainmodel.Login{}, nil},
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

func (ts *TestSuite) TestLogIn() {
	credentials := security.Credentials{}

	login := domainmodel.Login{}

	auth := domainmodel.Auth{}

	newAuth := domainmodel.Auth{}

	tokenExpTimeInSec := fake.Number(2, 10)

	token := ""

	errorType := customerror.NoType

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInLoggingIn",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainmodel.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainmodel.Auth{
					UserID: login.UserID,
				}

				id = uuid.NewV4()

				newAuth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				token = fake.Word()

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{domainmodel.Auth{}, nil},
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
				credentials = security.Credentials{}

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
					{domainmodel.Login{}, nil},
					{nil},
					{domainmodel.Auth{}, nil},
					{domainmodel.Auth{}, nil},
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
				credentials = securitypkgfactory.NewCredentials(nil)

				returnArgs = ReturnArgs{
					{nil},
					{domainmodel.Login{}, customerror.New("failed")},
					{nil},
					{domainmodel.Auth{}, nil},
					{domainmodel.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheUsernameIsNotRegistered",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				returnArgs = ReturnArgs{
					{nil},
					{domainmodel.Login{}, nil},
					{nil},
					{domainmodel.Auth{}, nil},
					{domainmodel.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NotFound
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenVerifyingThePasswords",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainmodel.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{customerror.New("failed")},
					{domainmodel.Auth{}, nil},
					{domainmodel.Auth{}, nil},
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

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainmodel.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainmodel.Auth{
					UserID: login.UserID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{domainmodel.Auth{}, customerror.New("failed")},
					{domainmodel.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheUserIDIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainmodel.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainmodel.Auth{
					UserID: login.UserID,
				}

				id = uuid.NewV4()

				newAuth = domainmodel.Auth{
					ID:     id,
					UserID: login.UserID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{auth, nil},
					{domainmodel.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingANewAuth",
			SetUp: func(t *testing.T) {
				credentials = securitypkgfactory.NewCredentials(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainmodel.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainmodel.Auth{
					UserID: login.UserID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{domainmodel.Auth{}, nil},
					{domainmodel.Auth{}, customerror.New("failed")},
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

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainmodel.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainmodel.Auth{
					UserID: login.UserID,
				}

				id = uuid.NewV4()

				newAuth = domainmodel.Auth{
					ID:     id,
					UserID: login.UserID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{domainmodel.Auth{}, nil},
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

			security := new(mocksecurity.Security)
			security.On("VerifyPasswords", login.Password, credentials.Password).Return(returnArgs[2]...)

			authDatastoreRepository := new(authdatastoremockrepository.Repository)
			authDatastoreRepository.On("GetByUserID", login.UserID.String()).Return(returnArgs[3]...)
			authDatastoreRepository.On("Create", auth).Return(returnArgs[4]...)

			authN := new(mockauth.Auth)
			authN.On("CreateToken", newAuth, tokenExpTimeInSec).Return(returnArgs[5]...)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, security, validator, tokenExpTimeInSec)

			returnedToken, err := authService.LogIn(credentials)

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

func (ts *TestSuite) TestRenewToken() {
	auth := domainmodel.Auth{}

	tokenExpTimeInSec := fake.Number(2, 10)

	token := ""

	errorType := customerror.NoType

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRenewingTheToken",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				token = fake.Word()

				returnArgs = ReturnArgs{
					{nil},
					{token, nil},
				}
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAToken",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				returnArgs = ReturnArgs{
					{nil},
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
			validator.On("Validate", auth).Return(returnArgs[0]...)

			authN := new(mockauth.Auth)
			authN.On("CreateToken", auth, tokenExpTimeInSec).Return(returnArgs[1]...)

			authDatastoreRepository := new(authdatastoremockrepository.Repository)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)

			loginDatastoreRepository := new(logindatastoremockrepository.Repository)

			security := new(mocksecurity.Security)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, security, validator, tokenExpTimeInSec)

			returnedToken, err := authService.RenewToken(auth)

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

func (ts *TestSuite) TestModifyPassword() {
	id := ""

	passwords := security.Passwords{}

	login := domainmodel.Login{}

	updatedLogin := domainmodel.Login{}

	errorType := customerror.NoType

	tokenExpTimeInSec := fake.Number(2, 10)

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInModifyingThePassword",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				passwords = securitypkgfactory.NewPasswords(nil)

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainmodel.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: passwords.CurrentPassword,
				}

				updatedLogin = login
				updatedLogin.Password = passwords.NewPassword

				newLogin := updatedLogin

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{nil},
					{newLogin, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheIDIsNotValid",
			SetUp: func(t *testing.T) {
				id = ""

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
					{nil},
					{domainmodel.Login{}, nil},
					{nil},
					{domainmodel.Login{}, nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheEvaluatedPasswordsValuesAreNotValid",
			SetUp: func(t *testing.T) {
				passwords = security.Passwords{}

				returnArgs = ReturnArgs{
					{nil},
					{customerror.New("failed")},
					{domainmodel.Login{}, nil},
					{nil},
					{domainmodel.Login{}, nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingALoginByUsername",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{domainmodel.Login{}, customerror.New("failed")},
					{nil},
					{domainmodel.Login{}, nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheIDIsNotRegistered",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{domainmodel.Login{}, nil},
					{nil},
					{domainmodel.Login{}, nil},
				}

				errorType = customerror.NotFound
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOfInvalidPasswordHappensWhenVerifyingThePasswords",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainmodel.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: currentPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{customerror.Unauthorized.New("the password is invalid")},
					{domainmodel.Login{}, nil},
				}

				errorType = customerror.Unauthorized
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnotherErrorHappensWhenVerifyingThePasswords",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainmodel.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: currentPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{customerror.New("failed")},
					{domainmodel.Login{}, nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheNewPasswordISTheSameAsTheOneCurrentlyRegistered",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := currentPassword

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainmodel.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: currentPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{nil},
					{domainmodel.Login{}, nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenUpdatingTheLogin",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainmodel.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: currentPassword,
				}

				updatedLogin = login
				updatedLogin.Password = newPassword

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{nil},
					{domainmodel.Login{}, customerror.New("failed")},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			validator := new(mockvalidator.Validator)
			validator.On("ValidateWithTags", id, "nonzero, uuid").Return(returnArgs[0]...)
			validator.On("Validate", passwords).Return(returnArgs[1]...)

			loginDatastoreRepository := new(logindatastoremockrepository.Repository)
			loginDatastoreRepository.On("GetByUserID", id).Return(returnArgs[2]...)

			security := new(mocksecurity.Security)
			security.On("VerifyPasswords", login.Password, passwords.CurrentPassword).Return(returnArgs[3]...)

			loginDatastoreRepository.On("Update", updatedLogin.ID.String(), updatedLogin).Return(returnArgs[4]...)

			authDatastoreRepository := new(authdatastoremockrepository.Repository)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)

			authN := new(mockauth.Auth)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, security, validator, tokenExpTimeInSec)

			err := authService.ModifyPassword(id, passwords)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
			}
		})
	}
}

func (ts *TestSuite) TestLogOut() {
	id := ""

	errorType := customerror.NoType

	tokenExpTimeInSec := fake.Number(2, 10)

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInLoggingOut",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()

				authID, err := uuid.FromString(id)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     authID,
					"userID": userID,
				}

				auth := domainfactorymodel.NewAuth(args)

				returnArgs = ReturnArgs{
					{nil},
					{auth, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheIDIsNotValid",
			SetUp: func(t *testing.T) {
				id = ""

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
					{domainmodel.Auth{}, nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenDeletingTheAuth",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()

				returnArgs = ReturnArgs{
					{nil},
					{domainmodel.Auth{}, customerror.New("failed")},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			validator := new(mockvalidator.Validator)
			validator.On("ValidateWithTags", id, "nonzero, uuid").Return(returnArgs[0]...)

			authDatastoreRepository := new(authdatastoremockrepository.Repository)
			authDatastoreRepository.On("Delete", id).Return(returnArgs[1]...)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)

			loginDatastoreRepository := new(logindatastoremockrepository.Repository)

			authN := new(mockauth.Auth)
			security := new(mocksecurity.Security)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, security, validator, tokenExpTimeInSec)

			err := authService.LogOut(id)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
			}
		})
	}
}

func (ts *TestSuite) TestWithDBTrx() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrx := &gorm.DB{}

	authDatastoreRepositoryWithDBTrx := &authdatastoremockrepository.Repository{}
	userDatastoreRepositoryWithDBTrx := &userdatastoremockrepository.Repository{}
	loginDatastoreRepositoryWithDBTrx := &logindatastoremockrepository.Repository{}

	tokenExpTimeInSec := fake.Number(2, 10)

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSettingRepositoriesWithDatabaseTransaction",
			SetUp: func(t *testing.T) {
				dbTrx = db

				authDatastoreRepositoryWithDBTrx = &authdatastoremockrepository.Repository{}
				userDatastoreRepositoryWithDBTrx = &userdatastoremockrepository.Repository{}
				loginDatastoreRepositoryWithDBTrx = &logindatastoremockrepository.Repository{}

				returnArgs = ReturnArgs{
					{authDatastoreRepositoryWithDBTrx},
					{userDatastoreRepositoryWithDBTrx},
					{loginDatastoreRepositoryWithDBTrx},
				}
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := new(authdatastoremockrepository.Repository)
			authDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[0]...)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)
			userDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[1]...)

			loginDatastoreRepository := new(logindatastoremockrepository.Repository)
			loginDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[2]...)

			authN := new(mockauth.Auth)
			security := new(mocksecurity.Security)
			validator := new(mockvalidator.Validator)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, security, validator, tokenExpTimeInSec)

			returnedAuthService := authService.WithDBTrx(dbTrx)

			if !tc.WantError {
				assert.NotEmpty(t, returnedAuthService, "Service interface is empty.")
				assert.Equal(t, authService, returnedAuthService, "Service interfaces are not the same.")
			}
		})
	}
}
