package login

import (
	"strings"

	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	logindatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/infrastructure/datastore/repository/login"
	persistententity "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/perentity"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

var initDB *gorm.DB

// New is the factory function that encapsulates the implementation related to login.
func New(db *gorm.DB) logindatastorerepository.IRepository {
	initDB = db
	return &Repository{
		DB: db,
	}
}

// Create is the function that creates a login in the database.
func (r *Repository) Create(login domainentity.Login) (domainentity.Login, error) {
	persistentLogin := persistententity.Login{}
	persistentLogin.FromDomain(login)

	if result := r.DB.Create(&persistentLogin); result.Error != nil {
		if strings.Contains(result.Error.Error(), "logins_user_id_key") {
			return domainentity.Login{}, customerror.Conflict.Newf("The user with id %s is already logged in", login.Username)
		}

		return domainentity.Login{}, result.Error
	}

	return persistentLogin.ToDomain(), nil
}

// GetByUsername is the function that gets a user by username from the database.
func (r *Repository) GetByUsername(username string) (domainentity.Login, error) {
	persistentLogin := persistententity.Login{}

	if result := r.DB.Find(&persistentLogin, "username=?", username); result.Error != nil {
		return domainentity.Login{}, result.Error
	}

	return persistentLogin.ToDomain(), nil
}

// GetByUsername is the function that gets a user by username from the database.
func (r *Repository) GetByUserID(userID string) (domainentity.Login, error) {
	persistentLogin := persistententity.Login{}

	if result := r.DB.Find(&persistentLogin, "user_id=?", userID); result.Error != nil {
		return domainentity.Login{}, result.Error
	}

	return persistentLogin.ToDomain(), nil
}

// Update is the function that updates a login by id in the database.
func (r *Repository) Update(id string, login domainentity.Login) (domainentity.Login, error) {
	persistentLogin := persistententity.Login{}
	persistentLogin.FromDomain(login)

	result := r.DB.Model(&persistentLogin).Where("id=?", id).Updates(&persistentLogin)
	if result.Error != nil {
		return domainentity.Login{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainentity.Login{}, customerror.NotFound.Newf("the login with id %s was not updated", id)
	}

	if result = r.DB.Find(&persistentLogin, "id=?", id); result.Error != nil {
		return domainentity.Login{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainentity.Login{}, customerror.NotFound.Newf("the login id %s was not found", id)
	}

	return persistentLogin.ToDomain(), nil
}

// Delete is the function that deletes a login by id from the database.
func (r *Repository) Delete(id string) (domainentity.Login, error) {
	persistentLogin := persistententity.Login{}

	result := r.DB.Find(&persistentLogin, "id=?", id)
	if result.Error != nil {
		return domainentity.Login{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainentity.Login{}, customerror.NotFound.Newf("the login with id %s was not found", id)
	}

	if result = r.DB.Delete(&persistentLogin); result.Error != nil {
		return domainentity.Login{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainentity.Login{}, customerror.NotFound.Newf("the login with id %s was not deleted", id)
	}

	return persistentLogin.ToDomain(), nil
}

// WithDBTrx is the function that enables the repository with database transaction.
func (r *Repository) WithDBTrx(dbTrx *gorm.DB) logindatastorerepository.IRepository {
	if dbTrx == nil {
		r.DB = initDB
		return r
	}

	r.DB = dbTrx
	return r
}
