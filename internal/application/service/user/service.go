package user

import (
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	userservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/service/user"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/infrastructure/datastore/repository/user"
	validatorpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/validator"
	"gorm.io/gorm"
)

type Service struct {
	UserDatastoreRepository userdatastorerepository.IRepository
	Validator               validatorpkg.IValidator
}

// New is the factory function that encapsulates the implementation related to user service.
func New(userDatastoreRepository userdatastorerepository.IRepository, validator validatorpkg.IValidator) userservice.IService {
	return &Service{
		UserDatastoreRepository: userDatastoreRepository,
	}
}

// GetAll is the function that deals with the user repository for getting all users.
func (u *Service) GetAll() (domainentity.Users, error) {
	users, err := u.UserDatastoreRepository.GetAll()
	if err != nil {
		return domainentity.Users{}, err
	}

	return users, nil
}

// WithDBTrx is the function that enables the service with database transaction.
func (u *Service) WithDBTrx(dbTrx *gorm.DB) userservice.IService {
	u.UserDatastoreRepository = u.UserDatastoreRepository.WithDBTrx(dbTrx)

	return u
}
