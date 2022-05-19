package auth

import (
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/security"
	"gorm.io/gorm"
)

// IService interface is a collection of function signatures that represents the auth's service contract.
type IService interface {
	Register(credentials securitypkg.Credentials) (string, error)
	LogIn(credentials securitypkg.Credentials) (string, error)
	RenewToken(auth domainmodel.Auth) (string, error)
	ModifyPassword(id string, passwords securitypkg.Passwords) error
	LogOut(id string) error
	WithDBTrx(dbTrx *gorm.DB) IService
}
