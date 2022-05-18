package user

import (
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	"gorm.io/gorm"
)

// IRepository interface is a collection of function signatures that represents the user's datastore repository contract.
type IRepository interface {
	Create(user domainmodel.User) (domainmodel.User, error)
	GetAll() (domainmodel.Users, error)
	WithDBTrx(dbTrx *gorm.DB) IRepository
}
