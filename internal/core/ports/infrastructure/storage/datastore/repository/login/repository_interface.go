package login

import (
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	"gorm.io/gorm"
)

// IRepository interface is a collection of function signatures that represents the login's repository contract.
type IRepository interface {
	Create(login domainmodel.Login) (domainmodel.Login, error)
	GetByUsername(username string) (domainmodel.Login, error)
	GetByUserID(userID string) (domainmodel.Login, error)
	Update(id string, login domainmodel.Login) (domainmodel.Login, error)
	Delete(id string) (domainmodel.Login, error)
	WithDBTrx(dbTrx *gorm.DB) IRepository
}
