package login

import (
	"strings"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/repository/login"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/model"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
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
func (r *Repository) Create(login domainmodel.Login) (domainmodel.Login, error) {
	loginDatastore := datastoremodel.Login{}
	loginDatastore.FromDomain(login)

	if result := r.DB.Create(&loginDatastore); result.Error != nil {
		if strings.Contains(result.Error.Error(), "logins_user_id_key") {
			return domainmodel.Login{}, customerror.Conflict.Newf("The user with id %s is already logged in", login.Username)
		}

		return domainmodel.Login{}, result.Error
	}

	return loginDatastore.ToDomain(), nil
}

// GetByUsername is the function that gets a user by username from the database.
func (r *Repository) GetByUsername(username string) (domainmodel.Login, error) {
	loginDatastore := datastoremodel.Login{}

	if result := r.DB.Find(&loginDatastore, "username=?", username); result.Error != nil {
		return domainmodel.Login{}, result.Error
	}

	return loginDatastore.ToDomain(), nil
}

// GetByUsername is the function that gets a user by username from the database.
func (r *Repository) GetByUserID(userID string) (domainmodel.Login, error) {
	loginDatastore := datastoremodel.Login{}

	if result := r.DB.Find(&loginDatastore, "user_id=?", userID); result.Error != nil {
		return domainmodel.Login{}, result.Error
	}

	return loginDatastore.ToDomain(), nil
}

// Update is the function that updates a login by id in the database.
func (r *Repository) Update(id string, login domainmodel.Login) (domainmodel.Login, error) {
	loginDatastore := datastoremodel.Login{}
	loginDatastore.FromDomain(login)

	result := r.DB.Model(&loginDatastore).Where("id=?", id).Updates(&loginDatastore)
	if result.Error != nil {
		return domainmodel.Login{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainmodel.Login{}, customerror.NotFound.Newf("the login with id %s was not updated", id)
	}

	if result = r.DB.Find(&loginDatastore, "id=?", id); result.Error != nil {
		return domainmodel.Login{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainmodel.Login{}, customerror.NotFound.Newf("the login id %s was not found", id)
	}

	return loginDatastore.ToDomain(), nil
}

// Delete is the function that deletes a login by id from the database.
func (r *Repository) Delete(id string) (domainmodel.Login, error) {
	loginDatastore := datastoremodel.Login{}

	result := r.DB.Find(&loginDatastore, "id=?", id)
	if result.Error != nil {
		return domainmodel.Login{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainmodel.Login{}, customerror.NotFound.Newf("the login with id %s was not found", id)
	}

	if result = r.DB.Delete(&loginDatastore); result.Error != nil {
		return domainmodel.Login{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainmodel.Login{}, customerror.NotFound.Newf("the login with id %s was not deleted", id)
	}

	return loginDatastore.ToDomain(), nil
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
