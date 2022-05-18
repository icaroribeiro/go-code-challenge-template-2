package user

import (
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	"gorm.io/gorm"
)

// IService interface is a collection of function signatures that represents the user's service contract.
type IService interface {
	GetAll() (domainmodel.Users, error)
	WithDBTrx(dbTrx *gorm.DB) IService
}
