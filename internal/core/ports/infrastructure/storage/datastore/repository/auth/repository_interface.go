package auth

import (
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	"gorm.io/gorm"
)

// IRepository interface is a collection of function signatures that represents the auth's datastore repository contract.
type IRepository interface {
	Create(auth domainmodel.Auth) (domainmodel.Auth, error)
	GetByUserID(userID string) (domainmodel.Auth, error)
	Delete(id string) (domainmodel.Auth, error)
	WithDBTrx(dbTrx *gorm.DB) IRepository
}
