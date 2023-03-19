package auth

import (
	"strings"

	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	authdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/infrastructure/datastore/repository/auth"
	persistententity "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/perentity"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

var initDB *gorm.DB

// New is the factory function that encapsulates the implementation related to auth.
func New(db *gorm.DB) authdatastorerepository.IRepository {
	initDB = db
	return &Repository{
		DB: db,
	}
}

// Create is the function that creates an auth in the datastore.
func (r *Repository) Create(auth domainentity.Auth) (domainentity.Auth, error) {
	persistentAuth := persistententity.Auth{}
	persistentAuth.FromDomain(auth)

	result := r.DB.Create(&persistentAuth)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "auths_user_id_key") {
			persistentLogin := persistententity.Login{}

			if result := r.DB.Find(&persistentLogin, "user_id=?", persistentAuth.UserID); result.Error != nil {
				return domainentity.Auth{}, result.Error
			}

			if result.RowsAffected == 0 && persistentLogin.IsEmpty() {
				return domainentity.Auth{}, customerror.NotFound.Newf("the login record with user id %s was not found", persistentAuth.UserID)
			}

			return domainentity.Auth{}, customerror.Conflict.Newf("The user with id %s is already logged in", persistentAuth.UserID)
		}

		return domainentity.Auth{}, result.Error
	}

	return persistentAuth.ToDomain(), nil
}

// GetByUserID is the function that gets an auth by user id from the datastore.
func (r *Repository) GetByUserID(userID string) (domainentity.Auth, error) {
	persistentAuth := persistententity.Auth{}

	if result := r.DB.Find(&persistentAuth, "user_id=?", userID); result.Error != nil {
		return domainentity.Auth{}, result.Error
	}

	return persistentAuth.ToDomain(), nil
}

// Delete is the function that deletes an auth by id from the datastore.
func (r *Repository) Delete(id string) (domainentity.Auth, error) {
	persistentAuth := persistententity.Auth{}

	result := r.DB.Find(&persistentAuth, "id=?", id)
	if result.Error != nil {
		return domainentity.Auth{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainentity.Auth{}, customerror.NotFound.Newf("the auth with id %s was not found", id)
	}

	if result = r.DB.Delete(&persistentAuth); result.Error != nil {
		return domainentity.Auth{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainentity.Auth{}, customerror.NotFound.Newf("the auth with id %s was not deleted", id)
	}

	return persistentAuth.ToDomain(), nil
}

// WithDBTrx is the function that enables the repository with datastore transaction.
func (r *Repository) WithDBTrx(dbTrx *gorm.DB) authdatastorerepository.IRepository {
	if dbTrx == nil {
		r.DB = initDB
		return r
	}

	r.DB = dbTrx
	return r
}
