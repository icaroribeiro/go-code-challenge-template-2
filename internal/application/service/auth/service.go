package auth

import (
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	authservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/service/auth"
	authdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/infrastructure/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/infrastructure/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/infrastructure/datastore/repository/user"
	authpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/auth"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/security"
	validatorpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/validator"
	"gorm.io/gorm"
)

type Service struct {
	AuthDatastoreRepository  authdatastorerepository.IRepository
	LoginDatastoreRepository logindatastorerepository.IRepository
	UserDatastoreRepository  userdatastorerepository.IRepository
	Validator                validatorpkg.IValidator
	AuthN                    authpkg.IAuth
	Security                 securitypkg.ISecurity
	TokenExpTimeInSec        int
}

// New is the factory function that encapsulates the implementation related to auth.
func New(authDatastoreRepository authdatastorerepository.IRepository,
	loginDatastoreRepository logindatastorerepository.IRepository,
	userDatastoreRepository userdatastorerepository.IRepository,
	authN authpkg.IAuth,
	security securitypkg.ISecurity,
	validator validatorpkg.IValidator,
	tokenExpTimeInSec int) authservice.IService {

	return &Service{
		AuthDatastoreRepository:  authDatastoreRepository,
		UserDatastoreRepository:  userDatastoreRepository,
		LoginDatastoreRepository: loginDatastoreRepository,
		AuthN:                    authN,
		Validator:                validator,
		Security:                 security,
		TokenExpTimeInSec:        tokenExpTimeInSec,
	}
}

// Register is the function that registers the user to the system.
func (a *Service) Register(credentials securitypkg.Credentials) (string, error) {
	if err := a.Validator.Validate(credentials); err != nil {
		return "", customerror.BadRequest.New(err.Error())
	}

	login, err := a.LoginDatastoreRepository.GetByUsername(credentials.Username)
	if err != nil {
		return "", err
	}

	if !login.IsEmpty() {
		return "", customerror.Conflict.Newf("the username %s is already registered", credentials.Username)
	}

	user := domainentity.User{
		Username: credentials.Username,
	}

	UserFactory, err := a.UserDatastoreRepository.Create(user)
	if err != nil {
		return "", err
	}

	login = domainentity.Login{
		UserID:   UserFactory.ID,
		Username: credentials.Username,
		Password: credentials.Password,
	}

	_, err = a.LoginDatastoreRepository.Create(login)
	if err != nil {
		return "", err
	}

	auth := domainentity.Auth{
		UserID: UserFactory.ID,
	}

	AuthFactory, err := a.AuthDatastoreRepository.Create(auth)
	if err != nil {
		return "", err
	}

	token, err := a.AuthN.CreateToken(AuthFactory, a.TokenExpTimeInSec)
	if err != nil {
		return "", err
	}

	return token, nil
}

// LogIn is the function that initializes the user access to the system.
func (a *Service) LogIn(credentials securitypkg.Credentials) (string, error) {
	if err := a.Validator.Validate(credentials); err != nil {
		return "", customerror.BadRequest.New(err.Error())
	}

	login, err := a.LoginDatastoreRepository.GetByUsername(credentials.Username)
	if err != nil {
		return "", err
	}

	if login.IsEmpty() {
		return "", customerror.NotFound.Newf("the username %s is not registered", credentials.Username)
	}

	if err = a.Security.VerifyPasswords(login.Password, credentials.Password); err != nil {
		return "", err
	}

	auth, err := a.AuthDatastoreRepository.GetByUserID(login.UserID.String())
	if err != nil {
		return "", err
	}

	if !auth.IsEmpty() {
		return "", customerror.Newf("the user with username %s is already logged in", credentials.Username)
	}

	auth = domainentity.Auth{
		UserID: login.UserID,
	}

	AuthFactory, err := a.AuthDatastoreRepository.Create(auth)
	if err != nil {
		return "", err
	}

	auth = AuthFactory

	token, err := a.AuthN.CreateToken(auth, a.TokenExpTimeInSec)
	if err != nil {
		return "", err
	}

	return token, nil
}

// RenewToken is the function that renews the token.
func (a *Service) RenewToken(auth domainentity.Auth) (string, error) {
	return a.AuthN.CreateToken(auth, a.TokenExpTimeInSec)
}

// ModifyPassword is the function that modifies the user's password.
func (a *Service) ModifyPassword(id string, passwords securitypkg.Passwords) error {
	if err := a.Validator.ValidateWithTags(id, "nonzero, uuid"); err != nil {
		return customerror.BadRequest.Newf("UserID: %s", err.Error())
	}

	if err := a.Validator.Validate(passwords); err != nil {
		return customerror.BadRequest.New(err.Error())
	}

	login, err := a.LoginDatastoreRepository.GetByUserID(id)
	if err != nil {
		return err
	}

	if login.IsEmpty() {
		return customerror.NotFound.New("the user who owns this token is not registered")
	}

	if err = a.Security.VerifyPasswords(login.Password, passwords.CurrentPassword); err != nil {
		if customerror.GetType(err) == customerror.Unauthorized {
			return customerror.Unauthorized.New("the current password did not match the one already registered")
		}

		return err
	}

	if passwords.NewPassword == passwords.CurrentPassword {
		return customerror.BadRequest.New("the new password is the same as the one currently registered")
	}

	login.Password = passwords.NewPassword

	_, err = a.LoginDatastoreRepository.Update(login.ID.String(), login)

	return err
}

// LogOut is the function that concludes the user access to the system.
func (a *Service) LogOut(id string) error {
	if err := a.Validator.ValidateWithTags(id, "nonzero, uuid"); err != nil {
		return customerror.BadRequest.New(err.Error())
	}

	_, err := a.AuthDatastoreRepository.Delete(id)

	return err
}

// WithDBTrx is the function that enables the service with database transaction.
func (a *Service) WithDBTrx(dbTrx *gorm.DB) authservice.IService {
	a.AuthDatastoreRepository = a.AuthDatastoreRepository.WithDBTrx(dbTrx)

	a.UserDatastoreRepository = a.UserDatastoreRepository.WithDBTrx(dbTrx)

	a.LoginDatastoreRepository = a.LoginDatastoreRepository.WithDBTrx(dbTrx)

	return a
}
