package login

import (
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	"gorm.io/gorm"
)

// IRepository interface is a collection of function signatures that represents the login's repository contract.
type IRepository interface {
	Create(login domainentity.Login) (domainentity.Login, error)
	GetByUsername(username string) (domainentity.Login, error)
	GetByUserID(userID string) (domainentity.Login, error)
	Update(id string, login domainentity.Login) (domainentity.Login, error)
	Delete(id string) (domainentity.Login, error)
	WithDBTrx(dbTrx *gorm.DB) IRepository
}
