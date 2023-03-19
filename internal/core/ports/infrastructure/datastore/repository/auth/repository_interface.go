package auth

import (
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	"gorm.io/gorm"
)

// IRepository interface is a collection of function signatures that represents the auth's datastore repository contract.
type IRepository interface {
	Create(auth domainentity.Auth) (domainentity.Auth, error)
	GetByUserID(userID string) (domainentity.Auth, error)
	Delete(id string) (domainentity.Auth, error)
	WithDBTrx(dbTrx *gorm.DB) IRepository
}
