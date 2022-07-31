package user

import (
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	"gorm.io/gorm"
)

// IRepository interface is a collection of function signatures that represents the user's datastore repository contract.
type IRepository interface {
	Create(user domainentity.User) (domainentity.User, error)
	GetAll() (domainentity.Users, error)
	WithDBTrx(dbTrx *gorm.DB) IRepository
}
