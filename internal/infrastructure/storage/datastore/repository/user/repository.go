package user

import (
	"strings"

	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/repository/user"
	datastoreentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/entity"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

var initDB *gorm.DB

// New is the factory function that encapsulates the implementation related to user repository.
func New(db *gorm.DB) userdatastorerepository.IRepository {
	initDB = db
	return &Repository{
		DB: db,
	}
}

// Create is the function that creates a user in the database.
func (r *Repository) Create(user domainentity.User) (domainentity.User, error) {
	userDatastore := datastoreentity.User{}
	userDatastore.FromDomain(user)

	if result := r.DB.Create(&userDatastore); result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			return domainentity.User{}, customerror.Conflict.New(result.Error.Error())
		}

		return domainentity.User{}, result.Error
	}

	return userDatastore.ToDomain(), nil
}

// GetAll is the function that gets the list of all users from the database.
func (r *Repository) GetAll() (domainentity.Users, error) {
	usersDatastore := datastoreentity.Users{}

	if result := r.DB.Find(&usersDatastore); result.Error != nil {
		return domainentity.Users{}, result.Error
	}

	return usersDatastore.ToDomain(), nil
}

// WithDBTrx is the function that enables the repository with database transaction.
func (r *Repository) WithDBTrx(dbTrx *gorm.DB) userdatastorerepository.IRepository {
	if dbTrx == nil {
		r.DB = initDB
		return r
	}

	r.DB = dbTrx
	return r
}
